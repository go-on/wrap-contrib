package wraps

import (
	"net/http"
	"github.com/go-on/wrap-contrib/helper"

	"github.com/go-on/wrap"
)

// ContentType writes the content type if the inner handler was successful
// and did not set a content-type
type ContentType string

// ServeHandle serves the given request with the inner handler and after that
// writes the content type, if the inner handler was successful
// and did not set a content-type
func (c ContentType) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	buf := helper.NewResponseBuffer(wr)
	inner.ServeHTTP(buf, req)
	if buf.IsOk() {
		wr.Header().Set("Content-Type", string(c))
	}
	buf.WriteTo(wr)
}

// Wrap wraps the given inner handler with the returned handler
func (c ContentType) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(c, inner)
}

// "application/json; charset=utf-8"
// text/plain; charset=utf-8
// text/css; charset=utf-8
// text/html; charset=utf-8
// application/javascript; charset=utf-8

var JSONContentType = ContentType("application/json; charset=utf-8")
var TextContentType = ContentType("text/plain; charset=utf-8")
var CSSContentType = ContentType("text/css; charset=utf-8")
var HTMLContentType = ContentType("text/html; charset=utf-8")
var JavaScriptContentType = ContentType("application/javascript; charset=utf-8")
