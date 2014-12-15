package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

// Catcher provides a Catch method that is called if a http.Handler recovered
// from a panic
type Catcher interface {
	// Catch is called if a http.Handler recovered from a panic
	// It is given the responsewriter and request and the return
	// value from the recover() call
	// The given ResponseWriter is only to enable Catch to write on it
	// and should not expected to have any bits written so far, Catch
	// should not examinate the ResponseWriter.
	Catch(recovered interface{}, w http.ResponseWriter, r *http.Request)
}

// CatchFunc is a function fullfilling the Catcher interface
type CatchFunc func(recovered interface{}, w http.ResponseWriter, r *http.Request)

// Catch fullfills the Catcher interface
func (c CatchFunc) Catch(recovered interface{}, w http.ResponseWriter, r *http.Request) {
	c(recovered, w, r)
}

// ServeHTTPNext serves the given request by letting the next serve a ResponseBuffer and
// catching any panics. If no panic happened, the ResponseBuffer is flushed to the ResponseWriter
// Otherwise the CatchFunc is called.
func (c CatchFunc) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	checked := wrap.NewPeek(wr, func(ck *wrap.Peek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})

	defer func() {
		if p := recover(); p != nil {
			c(p, wr, req)
		} else {
			checked.FlushMissing()
		}
	}()

	next.ServeHTTP(checked, req)
}

// Wrap implements the wrap.Wrapper interface
func (c CatchFunc) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(c).Wrap(next)
}

// Catch returns a CatchFunc for a Catcher
func Catch(c Catcher) wrap.Wrapper {
	return CatchFunc(c.Catch)
}
