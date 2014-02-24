package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
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
	buf := helper.NewResponseBuffer(wr)
	if req.Method == "HEAD" {
		req.Method = "GET"

		defer func() {
			req.Method = "HEAD"
		}()

		inner.ServeHTTP(buf, req)
		if buf.HasChanged() {
			buf.WriteHeadersTo(wr)
			if buf.Code != 0 {
				wr.WriteHeader(buf.Code)
				return
			}
			wr.WriteHeader(200)
			return
		}
		return
	}
	inner.ServeHTTP(wr, req)
}
