package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

// AfterFunc is of the type http.HandlerFunc
// and provides a wrap.Wrapper that calls itself after
// the inner handler has been called
type AfterFunc func(http.ResponseWriter, *http.Request)

// ServeHandle serves the given request with the inner handler and after that
// with the AfterFunc
func (a AfterFunc) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	inner.ServeHTTP(wr, req)
	a(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (a AfterFunc) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(a, inner)
}

// After returns an AfterFunc for a http.Handler
func After(h http.Handler) wrap.Wrapper {
	return AfterFunc(h.ServeHTTP)
}
