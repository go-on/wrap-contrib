package wraps

import (
	"net/http"
	"strings"

	"github.com/go-on/wrap"
)

// RemoveResponseHeader removes response headers that are identical to the string
// or have if as prefix
type RemoveResponseHeader string

// ServeHandle removes the response headers that are identical to the string
// or have if as prefix after the next handler is run
func (rh RemoveResponseHeader) ServeHTTPNext(next http.Handler, w http.ResponseWriter, r *http.Request) {
	bodyWritten := false
	comp := strings.TrimSpace(strings.ToLower(string(rh)))

	checked := wrap.NewPeek(w, func(ck *wrap.Peek) bool {
		hd := ck.Header()
		for k := range hd {
			k = strings.TrimSpace(strings.ToLower(k))
			if strings.HasPrefix(k, comp) {
				hd.Del(k)
			}
		}
		ck.FlushHeaders()
		ck.FlushCode()
		bodyWritten = true
		return true
	})

	next.ServeHTTP(checked, r)

	if !bodyWritten {
		hd := checked.Header()
		for k := range hd {
			k = strings.TrimSpace(strings.ToLower(k))
			if strings.HasPrefix(k, comp) {
				hd.Del(k)
			}
		}
	}
	checked.FlushMissing()
}

// Wrap wraps the given next handler with the returned handler
func (rh RemoveResponseHeader) Wrap(in http.Handler) http.Handler {
	return wrap.NextHandler(rh).Wrap(in)
}
