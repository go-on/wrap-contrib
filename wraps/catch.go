package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
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

// ServeHandle serves the given request by letting the next serve a ResponseBuffer and
// catching any panics. If no panic happened, the ResponseBuffer is flushed to the ResponseWriter
// Otherwise the CatchFunc is called.
func (c CatchFunc) ServeHandle(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			c(p, wr, req)
		}
	}()

	checked := wrap.NewRWPeek(wr, func(ck *wrap.RWPeek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})

	next.ServeHTTP(checked, req)

	checked.FlushMissing()

}

// Wrap wraps the given next handler with the returned handler
func (c CatchFunc) Wrap(next http.Handler) http.Handler {
	return wrap.ServeHandle(c, next)
}

// Catch returns a CatchFunc for a Catcher
func Catch(c Catcher) wrap.Wrapper {
	return CatchFunc(c.Catch)
}
