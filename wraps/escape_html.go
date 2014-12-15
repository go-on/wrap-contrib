package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

type escapeHTML struct{}

// Wrap wraps the given next handler within a http.HandlerFunc that
// calls the next handlers ServeHTTP method with an EscapeHTMLResponseWriter
func (e escapeHTML) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&wrap.EscapeHTML{wr}, req)
	})
}

func (e escapeHTML) WrapFunc(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&wrap.EscapeHTML{wr}, req)
	})
}

// EscapeHTML wraps the next handler by replacing the response writer with an EscapeHTMLResponseWriter
// that escapes html special chars while writing to the underlying response writer
var EscapeHTML = escapeHTML{}
