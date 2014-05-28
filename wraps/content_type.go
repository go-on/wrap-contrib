package wraps

import (
	"github.com/go-on/wrap-contrib/helper"
	"net/http"

	"github.com/go-on/wrap"
)

// ContentType writes the content type if the inner handler was successful
// and did not set a content-type
type ContentType string

// SetContentType sets the content type in the given ResponseWriter
func (c ContentType) SetContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", string(c))
}

// ServeHandle serves the given request with the inner handler and after that
// writes the content type, if the inner handler was successful
// and did not set a content-type
func (c ContentType) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	buf := helper.NewResponseBuffer(wr)
	inner.ServeHTTP(buf, req)
	if buf.IsOk() {
		c.SetContentType(wr)
	}
	buf.WriteAllTo(wr)
}

// Wrap wraps the given inner handler with the returned handler
func (c ContentType) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(c, inner)
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
