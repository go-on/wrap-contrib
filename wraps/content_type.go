package wraps

import (
	"fmt"
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

// ContentType writes the content type if the next handler was successful
// and did not set a content-type
type ContentType string

// SetContentType sets the content type in the given ResponseWriter
func (c ContentType) SetContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", string(c))
}

func (c ContentType) String() string {
	return fmt.Sprintf("<ContentType %#v>", string(c))
}

// ServeHandle serves the given request with the next handler and after that
// writes the content type, if the next handler was successful
// and did not set a content-type
func (c ContentType) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {

	var bodyWritten = false
	checked := wrap.NewPeek(wr, func(ck *wrap.Peek) bool {
		if ck.IsOk() {
			c.SetContentType(ck)
		}
		ck.FlushHeaders()
		ck.FlushCode()
		bodyWritten = true
		return true
	})

	next.ServeHTTP(checked, req)

	if !bodyWritten {
		c.SetContentType(checked)
	}
	checked.FlushMissing()
}

// Wrap wraps the given next handler with the returned handler
func (c ContentType) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(c).Wrap(next)
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
