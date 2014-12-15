package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

type setResponseHeader struct {
	Key, Val string
}

// SetResponseHeader sets a response header
func SetResponseHeader(key, val string) wrap.Wrapper {
	return &setResponseHeader{key, val}
}

// ServeHTTPNext sets the header of the key to the value and calls
// the next handler after that
func (s *setResponseHeader) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set(s.Key, s.Val)
	next.ServeHTTP(wr, req)
}

// Wrap wraps the given next handler with the returned handler
func (s *setResponseHeader) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(s).Wrap(next)
}
