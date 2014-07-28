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

func (f *fallback) ServeHandle(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	bodywritten := false
	checked := helper.NewCheckedResponseWriter(wr, func(ck *helper.CheckedResponseWriter) bool {
		if _, has := f.ignoreCodes[ck.Code]; !has {
			ck.WriteHeadersTo(wr)
			ck.WriteCodeTo(wr)
			bodywritten = true
			return true
		}
		ck.Reset()
		return false
	})

	for _, h := range f.handlers {
		h.ServeHTTP(checked, req)
		if checked.HasChanged() {
			if _, has := f.ignoreCodes[checked.Code]; !has {
				if !bodywritten {
					checked.WriteHeadersTo(wr)
					checked.WriteCodeTo(wr)
				}
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
