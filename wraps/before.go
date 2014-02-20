package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

// BeforeFunc is of the type http.HandlerFunc
// and provides a wrap.Wrapper that calls itself before
// the inner handler has been called
type BeforeFunc func(http.ResponseWriter, *http.Request)

// ServeHandle serves the given request with the BeforeFunc and after that
// with the inner handler
func (a BeforeFunc) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	a(wr, req)
	inner.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (a BeforeFunc) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(a, inner)
}

// Before returns an BeforeFunc for a http.Handler
func Before(h http.Handler) wrap.Wrapper {
	return BeforeFunc(h.ServeHTTP)
}
