package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

// GuardFunc is a wrap.Wapper and http.HandlerFunc that may operate on the ResponseWriter
// If it does so, the wrapper prevents the next Handler from serving.
type GuardFunc func(http.ResponseWriter, *http.Request)

// ServeHTTPNext lets the GuardFunc serve to a ResponseBuffer and if it changed something
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

// Wrap implements the wrap.Wrapper interface
func (g GuardFunc) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(g).Wrap(next)
}

// Guard returns a GuardFunc for a http.Handler
func Guard(h http.Handler) wrap.Wrapper {
	return GuardFunc(h.ServeHTTP)
}
