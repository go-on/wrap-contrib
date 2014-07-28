package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type first []http.Handler

func (f first) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	checked := wrap.NewRWPeek(wr, func(ck *wrap.RWPeek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})

	for _, h := range f {
		h.ServeHTTP(checked, req)
		if checked.HasChanged() {
			checked.FlushMissing()
			return
		}
	}
	inner.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (f first) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(f, inner)
}

// First will try all given handler until
// the first one returns something
func First(handler ...http.Handler) wrap.Wrapper { return first(handler) }

// FirstFunc is like First but for http.HandlerFuncs
func FirstFunc(handlerFn ...func(w http.ResponseWriter, r *http.Request)) wrap.Wrapper {
	h := make([]http.Handler, len(handlerFn))
	for i, fn := range handlerFn {
		h[i] = http.HandlerFunc(fn)
	}
	return first(h)
}
