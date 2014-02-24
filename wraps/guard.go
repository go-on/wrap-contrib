package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
)

// GuardFunc is a wrap.Wapper and http.HandlerFunc that may operate on the ResponseWriter
// If it does so, the wrapper prevents the inner Handler from serving.
type GuardFunc func(http.ResponseWriter, *http.Request)

// ServeHandle lets the GuardFunc serve to a ResponseBuffer and if it changed something
// the Response is send to the ResponseWriter, preventing the inner http.Handler from
// executing. Otherwise the inner handler serves the request.
func (g GuardFunc) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	buf := helper.NewResponseBuffer(wr)
	g(buf, req)
	if buf.HasChanged() {
		buf.WriteTo(wr)
		return
	}
	inner.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (g GuardFunc) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(g, inner)
}

// Guard returns a GuardFunc for a http.Handler
func Guard(h http.Handler) wrap.Wrapper {
	return GuardFunc(h.ServeHTTP)
}
