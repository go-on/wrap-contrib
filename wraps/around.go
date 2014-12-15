package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

type around struct{ before, after http.Handler }

func (a around) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	a.before.ServeHTTP(wr, req)
	next.ServeHTTP(wr, req)
	a.after.ServeHTTP(wr, req)
}

// Wrap implements the wrap.Wrapper interface
func (a around) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(a).Wrap(next)
}

// Around returns a wrapper that calls the given before and after handler
// before and after the next handler when serving
func Around(before, after http.Handler) wrap.Wrapper {
	return around{before, after}
}

// AroundFunc returns a wrapper that acts like Around, but for HandleFuncs
func AroundFunc(before, after func(http.ResponseWriter, *http.Request)) wrap.Wrapper {
	return around{http.HandlerFunc(before), http.HandlerFunc(after)}
}
