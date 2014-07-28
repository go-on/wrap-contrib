package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type head struct{}

// Head will check if the request method is HEAD.
// if so, it will change the method to GET, call the handler
// with a ResponseBuffer and return only the header information
// to the client.
//
// For non HEAD methods, it simply pass the request handling
// through to the http.handler
var _head = head{}

func Head() wrap.Wrapper {
	return _head
}

// Wrap wraps the given inner handler with the returned handler
func (h head) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(h, inner)
}

// ServeHandle handles serves the request by transforming a HEAD request to a GET request
// for the inner handler and then remove the body from the response.
// Non HEAD reqeusts are not affected
func (h head) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	if req.Method == "HEAD" {
		req.Method = "GET"

		checked := wrap.NewRWPeek(wr, func(ck *wrap.RWPeek) bool {
			ck.FlushHeaders()
			ck.FlushCode()
			return false // write no body
		})

		defer func() {
			req.Method = "HEAD"
		}()

		inner.ServeHTTP(checked, req)

		checked.FlushMissing()
		return
	}
	inner.ServeHTTP(wr, req)
}
