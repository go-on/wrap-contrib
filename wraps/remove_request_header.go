package wraps

import (
	"net/http"
	"strings"

	"github.com/go-on/wrap"
)

// RemoveRequestHeader removes request headers that are identical to the string
// or have it as prefix
type RemoveRequestHeader string

// ServeHandle removes request headers that are identical to the string
// or have it as prefix. Then the inner http.Handler is called
func (rh RemoveRequestHeader) ServeHandle(inner http.Handler, w http.ResponseWriter, r *http.Request) {
	comp := strings.TrimSpace(strings.ToLower(string(rh)))
	for k := range r.Header {
		k = strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(k, comp) {
			r.Header.Del(k)
		}
	}
	inner.ServeHTTP(w, r)
}

// Wrap wraps the given inner handler with the returned handler
func (rh RemoveRequestHeader) Wrap(in http.Handler) http.Handler {
	return wrap.ServeHandle(rh, in)
}
