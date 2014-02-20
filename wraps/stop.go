package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type stop struct{}

// ServeHTTP does nothing
func (stop) ServeHTTP(wr http.ResponseWriter, req *http.Request) {}

// Wrap returns a handler that does do nothing
func (s stop) Wrap(inner http.Handler) http.Handler {
	return s
}

// Stop is a wrapper that does no processing but simply prevents further execution of
// inner wrappers
var _stop = stop{}

func Stop() wrap.Wrapper {
	return _stop
}
