package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

// GuardFunc is a wrap.Wapper and http.HandlerFunc that may operate on the ResponseWriter
// If it does so, the wrapper prevents the next Handler from serving.
type GuardFunc func(http.ResponseWriter, *http.Request)

// ServeHandle lets the GuardFunc serve to a ResponseBuffer and if it changed something
// the Response is send to the ResponseWriter, preventing the next http.Handler from
// executing. Otherwise the next handler serves the request.
func (g GuardFunc) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	checked := wrap.NewPeek(wr, func(ck *wrap.Peek) bool {
		ck.FlushHeaders()
		ck.FlushCode()
		return true
	})
	g(checked, req)
	if checked.HasChanged() {
		checked.FlushMissing()
		return
	}
	next.ServeHTTP(wr, req)
}

// Wrap wraps the given next handler with the returned handler
func (g GuardFunc) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(g).Wrap(next)
}

// Guard returns a GuardFunc for a http.Handler
func Guard(h http.Handler) wrap.Wrapper {
	return GuardFunc(h.ServeHTTP)
}
