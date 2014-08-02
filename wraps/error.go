package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

// Error a type based on error that should be saved by a wrap.Contexter (response writer)
type Error error

type errorHandler struct {
	http.Handler
}

func (e *errorHandler) Wrap(next http.Handler) http.Handler {

	// returns true, if error happened and was handled, otherwise false
	var handleError = func(rw http.ResponseWriter, req *http.Request) bool {
		var err error
		rw.(wrap.Contexter).Context(&err)
		if err != nil {
			e.ServeHTTP(rw, req)
			return true
		}
		// next.ServeHTTP(rw, req)
		return false
	}
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		bodywritten := false
		checked := wrap.NewPeek(rw, func(ck *wrap.Peek) bool {
			bodywritten = true
			if handleError(rw, req) {
				return false
			}
			ck.FlushMissing()
			return true
		})

		next.ServeHTTP(checked, req)

		if !bodywritten && !handleError(rw, req) {
			checked.FlushMissing()
		}
	}
	return f
}

// ErrorHandler returns a wrapper that requires the response writer to implement the
// wrap.Contexter interface and to support the *error type.
// When serving, it makes a fake run on the next handler and checks, if there is an error context
// and if so, it runs the given handler. Otherwise the writes of the next handler are flushed to the
// response writer
func ErrorHandler(h http.Handler) wrap.Wrapper {
	return &errorHandler{h}
}
