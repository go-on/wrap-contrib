package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type setRequestHeader struct {
	Key, Val string
}

func SetRequestHeader(key, val string) wrap.Wrapper {
	return &setRequestHeader{key, val}
}

// ServeHandle sets the header of the key to the value and calls
// the inner handler after that
func (s *setRequestHeader) ServeHTTPNext(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	req.Header.Set(s.Key, s.Val)
	inner.ServeHTTP(wr, req)
}

// Wrap wraps the given inner handler with the returned handler
func (s *setRequestHeader) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(s).Wrap(next)
}
