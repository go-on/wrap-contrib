package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type setResponseHeader struct {
	Key, Val string
}

func SetResponseHeader(key, val string) wrap.Wrapper {
	return &setResponseHeader{key, val}
}

// ServeHandle sets the header of the key to the value and calls
// the inner handler after that
func (s *setResponseHeader) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set(s.Key, s.Val)
	inner.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (s *setResponseHeader) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(s, inner)
}
