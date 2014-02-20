package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type around struct{ before, after http.Handler }

// ServeHandle serves the given request by calling the before handler followed by the inner
// handler and followed by the after handler
func (a around) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	a.before.ServeHTTP(wr, req)
	inner.ServeHTTP(wr, req)
	a.after.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (a around) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(a, inner)
}

// Around returns a wrapper that calls the given before and after handler
// before and after the inner handler when serving
func Around(before, after http.Handler) wrap.Wrapper {
	return around{before, after}
}

// AroundFunc returns a wrapper that acts like Around, but for HandleFuncs
func AroundFunc(before, after func(http.ResponseWriter, *http.Request)) wrap.Wrapper {
	return around{http.HandlerFunc(before), http.HandlerFunc(after)}
}
