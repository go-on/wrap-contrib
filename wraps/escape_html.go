package wraps

import (
	"net/http"

	"github.com/go-on/wrap-contrib/helper"
)

type escapeHTML struct{}

// Wrap wraps the given next handler within a http.HandlerFunc that
// calls the inner handlers ServeHTTP method with an EscapeHTMLResponseWriter
func (a escapeHTML) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&helper.EscapeHTMLResponseWriter{wr}, req)
	})
}

// EscapeHTML wraps the next handler by replacing the response writer with an EscapeHTMLResponseWriter
// that escapes html special chars while writing to the underlying response writer
var EscapeHTML = escapeHTML{}
