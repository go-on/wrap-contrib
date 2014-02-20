package wraps

import (
	"net/http"
	"strings"
	"github.com/go-on/wrap-contrib/helper"

	"github.com/go-on/wrap"
)

// RemoveResponseHeader removes response headers that are identical to the string
// or have if as prefix
type RemoveResponseHeader string

// ServeHandle removes the response headers that are identical to the string
// or have if as prefix after the inner handler is run
func (rh RemoveResponseHeader) ServeHandle(inner http.Handler, w http.ResponseWriter, r *http.Request) {
	buf := helper.NewResponseBuffer()
	inner.ServeHTTP(buf, r)

	comp := strings.TrimSpace(strings.ToLower(string(rh)))
	hd := buf.Header()
	for k := range hd {
		k = strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(k, comp) {
			hd.Del(k)
		}
	}

	buf.WriteTo(w)
}

// Wrap wraps the given inner handler with the returned handler
func (rh RemoveResponseHeader) Wrap(in http.Handler) http.Handler {
	return wrap.ServeHandle(rh, in)
}
