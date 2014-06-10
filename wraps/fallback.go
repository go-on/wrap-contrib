package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
)

type fallback struct {
	handlers    []http.Handler
	ignoreCodes map[int]struct{}
}

func (f *fallback) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	buf := helper.NewResponseBuffer(wr)
	for _, h := range f.handlers {
		h.ServeHTTP(buf, req)
		if buf.HasChanged() {
			if _, has := f.ignoreCodes[buf.Code]; !has {
				buf.WriteAllTo(wr)
				return
			}
			// remove any headers or status codes
			buf.Reset()
		}
	}
	inner.ServeHTTP(wr, req)
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
