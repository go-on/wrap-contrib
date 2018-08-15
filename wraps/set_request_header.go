package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type setRequestHeader struct {
	Key, Val string
}

// SetRequestHeader sets a request header
func SetRequestHeader(key, val string) wrap.Wrapper {
	return &setRequestHeader{key, val}
}

func (s *setRequestHeader) ServeHTTPNext(next http.Handler, wr http.ResponseWriter, req *http.Request) {
	req.Header.Set(s.Key, s.Val)
	next.ServeHTTP(wr, req)
}

// Wrap implements the wrapper interface
func (s *setRequestHeader) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(s).Wrap(next)
}
