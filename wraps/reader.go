package wraps

import (
	"io"
	"net/http"

	"github.com/go-on/wrap"
)

type reader struct {
	io.ReadSeeker
}

// Reader provides a wrap.Wrapper for a io.ReadSeeker
func Reader(r io.ReadSeeker) wrap.Wrapper {
	return &reader{r}
}

// ServeHTTP copies from the ReadSeeker starting at pos 0 to the ResponseWriter
// Any error results in Status Code 500
func (r *reader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, err := r.Seek(0, 0)
	if err != nil {
		rw.WriteHeader(500)
	}
	_, err = io.Copy(rw, r.ReadSeeker)
	if err != nil {
		rw.WriteHeader(500)
	}
}

// Wrap wraps the given inner handler with the returned handler
func (r *reader) Wrap(http.Handler) http.Handler { return r }
