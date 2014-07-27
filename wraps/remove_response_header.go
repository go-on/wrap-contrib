package wraps

import (
	"net/http"
	"strings"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
)

// RemoveResponseHeader removes response headers that are identical to the string
// or have if as prefix
type RemoveResponseHeader string

// ServeHandle removes the response headers that are identical to the string
// or have if as prefix after the next handler is run
func (rh RemoveResponseHeader) ServeHandle(next http.Handler, w http.ResponseWriter, r *http.Request) {
	checked := helper.NewCheckedResponseWriter(w, func(ck *helper.CheckedResponseWriter) bool {
		comp := strings.TrimSpace(strings.ToLower(string(rh)))
		hd := ck.Header()
		for k := range hd {
			k = strings.TrimSpace(strings.ToLower(k))
			if strings.HasPrefix(k, comp) {
				hd.Del(k)
			}
		}
		ck.WriteHeadersTo(w)
		ck.WriteCodeTo(w)
		return true
	})

	next.ServeHTTP(checked, r)
}

// Wrap wraps the given next handler with the returned handler
func (rh RemoveResponseHeader) Wrap(in http.Handler) http.Handler {
	return wrap.ServeHandle(rh, in)
}
