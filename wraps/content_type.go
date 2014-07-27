package wraps

import (
	"net/http"

	"github.com/go-on/wrap-contrib/helper"

	"github.com/go-on/wrap"
)

// ContentType writes the content type if the next handler was successful
// and did not set a content-type
type ContentType string

// SetContentType sets the content type in the given ResponseWriter
func (c ContentType) SetContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", string(c))
}

// ServeHandle serves the given request with the next handler and after that
// writes the content type, if the next handler was successful
// and did not set a content-type
func (c ContentType) ServeHandle(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	checked := helper.NewCheckedResponseWriter(wr, func(ck *helper.CheckedResponseWriter) bool {
		if ck.IsOk() {
			c.SetContentType(ck)
		}
		ck.WriteHeadersTo(wr)
		ck.WriteCodeTo(wr)
		return true
	})

	next.ServeHTTP(checked, req)
}

// Wrap wraps the given next handler with the returned handler
func (c ContentType) Wrap(next http.Handler) http.Handler {
	return wrap.ServeHandle(c, next)
}

var (
	JSONContentType       = ContentType("application/json; charset=utf-8")
	TextContentType       = ContentType("text/plain; charset=utf-8")
	CSSContentType        = ContentType("text/css; charset=utf-8")
	HTMLContentType       = ContentType("text/html; charset=utf-8")
	JavaScriptContentType = ContentType("application/javascript; charset=utf-8")
	RSSFeedContentType    = ContentType("application/rss+xml; charset=utf-8")
	AtomFeedContentType   = ContentType("application/atom+xml; charset=utf-8")
)
