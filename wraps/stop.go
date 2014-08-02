package wraps

import "net/http"

type stop struct{}

// ServeHTTP does nothing
func (stop) ServeHTTP(wr http.ResponseWriter, req *http.Request) {}

// Wrap returns a handler that does do nothing and does not run the next handler
// stopping the stack chain
func (s stop) Wrap(next http.Handler) http.Handler {
	return s
}

// Stop is a wrapper that does no processing but simply prevents further execution of
// next wrappers
var Stop = stop{}
