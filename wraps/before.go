package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

// BeforeFunc is of the type http.HandlerFunc
// and provides a wrap.Wrapper that calls itself before
// the next handler has been called
type BeforeFunc func(http.ResponseWriter, *http.Request)

func (b BeforeFunc) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	b(wr, req)
	next.ServeHTTP(wr, req)
}

// Wrap implements the wrap.Wrapper interface
func (b BeforeFunc) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(b).Wrap(next)
}

// Before returns an BeforeFunc for a http.Handler
func Before(h http.Handler) wrap.Wrapper {
	return BeforeFunc(h.ServeHTTP)
}
