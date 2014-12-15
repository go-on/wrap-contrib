package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

type head struct{}

var _head = head{}

// Head will check if the request method is HEAD.
// if so, it will change the method to GET, call the handler
// with a ResponseBuffer and return only the header information
// to the client.
//
// For non HEAD methods, it simply pass the request handling
// through to the http.handler
func Head() wrap.Wrapper {
	return _head
}

// Wrap wraps the given next handler with the returned handler
func (h head) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(h).Wrap(next)
}

// ServeHTTPNext handles serves the request by transforming a HEAD request to a GET request
// for the next handler and then remove the body from the response.
// Non HEAD requests are not affected
func (h head) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	if req.Method == "HEAD" {
		req.Method = "GET"

		checked := wrap.NewPeek(wr, func(ck *wrap.Peek) bool {
			ck.FlushHeaders()
			ck.FlushCode()
			return false // write no body
		})

		defer func() {
			req.Method = "HEAD"
		}()

		next.ServeHTTP(checked, req)

		checked.FlushMissing()
		return
	}
	next.ServeHTTP(wr, req)
}
