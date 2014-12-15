package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

// AfterFunc is of the type http.HandlerFunc
// and provides a wrap.Wrapper that calls itself after
// the inner handler has been called
type AfterFunc func(http.ResponseWriter, *http.Request)

func (a AfterFunc) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	next.ServeHTTP(wr, req)
	a(wr, req)
}

// Wrap implements the wrap.Wrapper interface
func (a AfterFunc) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(a).Wrap(next)
}

// After returns an AfterFunc for a http.Handler
func After(h http.Handler) wrap.Wrapper {
	return AfterFunc(h.ServeHTTP)
}
