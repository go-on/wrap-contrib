package wraps

import (
	"io"
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

type readseeker struct {
	io.ReadSeeker
}

// ReadSeeker provides a wrap.Wrapper for a io.ReadSeeker
func ReadSeeker(r io.ReadSeeker) wrap.Wrapper {
	return &readseeker{r}
}

// ServeHTTP copies from the ReadSeeker starting at pos 0 to the ResponseWriter
// Any error results in Status Code 500
func (r *readseeker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
func (r *readseeker) Wrap(http.Handler) http.Handler { return r }
