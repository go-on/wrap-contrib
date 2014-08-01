package wraps

import (
	"net/http"
	"strings"

	"github.com/go-on/wrap"
)

// RemoveRequestHeader removes request headers that are identical to the string
// or have it as prefix
type RemoveRequestHeader string

func (rh RemoveRequestHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	comp := strings.TrimSpace(strings.ToLower(string(rh)))
	for k := range r.Header {
		k = strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(k, comp) {
			r.Header.Del(k)
		}
	}
}

// ServeHandle removes request headers that are identical to the string
// or have it as prefix. Then the inner http.Handler is called
func (rh RemoveRequestHeader) ServeHTTPNext(inner http.Handler, w http.ResponseWriter, r *http.Request) {
	rh.ServeHTTP(w, r)
	inner.ServeHTTP(w, r)
}

// Wrap wraps the given inner handler with the returned handler
func (rh RemoveRequestHeader) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(rh).Wrap(next)
}
