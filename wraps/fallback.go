package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type fallback struct {
	handlers    []http.Handler
	ignoreCodes map[int]struct{}
}

func (f *fallback) ServeHandle(next http.Handler, wr http.ResponseWriter, req *http.Request) {

	checked := wrap.NewRWPeek(wr, func(ck *wrap.RWPeek) bool {
		if _, has := f.ignoreCodes[ck.Code]; !has {
			ck.FlushHeaders()
			ck.FlushCode()
			return true
		}
		ck.Reset()
		return false
	})

	for _, h := range f.handlers {
		h.ServeHTTP(checked, req)
		if checked.HasChanged() {
			if _, has := f.ignoreCodes[checked.Code]; !has {
				checked.FlushMissing()
				return
			}
		}
		checked.Reset()
	}
	next.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (f *fallback) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(f, inner)
}

// Fallback will try all given handler until
// the first one writes to the ResponseWriter body.
// It is similar to First, but ignores writes that did set one of the ignore status codes
func Fallback(ignoreCodes []int, handler ...http.Handler) wrap.Wrapper {
	fb := &fallback{handlers: handler, ignoreCodes: map[int]struct{}{}}
	for _, code := range ignoreCodes {
		fb.ignoreCodes[code] = struct{}{}
	}
	return fb
}

// FallbackFunc is like Fallback but for http.HandlerFuncs
func FallbackFunc(ignoreCodes []int, handlerFn ...func(w http.ResponseWriter, r *http.Request)) wrap.Wrapper {
	h := make([]http.Handler, len(handlerFn))
	for i, fn := range handlerFn {
		h[i] = http.HandlerFunc(fn)
	}
	return Fallback(ignoreCodes, h...)
}
