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
func (a BeforeFunc) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	a(wr, req)
	next.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (a BeforeFunc) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(a).Wrap(next)
}

// Before returns an BeforeFunc for a http.Handler
func Before(h http.Handler) wrap.Wrapper {
	return BeforeFunc(h.ServeHTTP)
}
