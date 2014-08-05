package wraphttpauth_test

import (
	"fmt"
	"net/http"

	"github.com/go-on/wrap"

	"github.com/abbot/go-http-auth"
)

// context is an example how a wrap.Contexter can be build in order to get used
// with wraphttpauth
type context struct {
	http.ResponseWriter
	authReq *auth.AuthenticatedRequest
}

// make sure to fulfill the Contexter interface
var _ wrap.Contexter = &context{}

// Context is an implementation for the Contexter interface.
//
// It receives a pointer to a value that is already stored inside the context.
// Values are distiguished by their type.
// Context sets the value of the given pointer to the value of the same type
// that is stored inside of the context.
// A pointer type that is not supported results in a panic.
func (c *context) Context(ctxPtr interface{}) (found bool) {
	found = true
	switch ty := ctxPtr.(type) {
	case *auth.AuthenticatedRequest:
		if c.authReq == nil {
			return false
		}
		*ty = *c.authReq
	default:
		panic(fmt.Sprintf("unsupported context: %T", ctxPtr))
	}
	return
}

// SetContext is an implementation for the Contexter interface.
//
// It receives a pointer to a value that will be stored inside the context.
// Values are distiguished by their type, that means that SetContext replaces
// and stored value of the same type.
// A pointer type that is not supported results in a panic.
func (c *context) SetContext(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *auth.AuthenticatedRequest:
		c.authReq = ty
	default:
		panic(fmt.Sprintf("unsupported context: %T", ctxPtr))
	}
}

// Wrap implements the wrap.Wrapper interface.
//
// When the request is served, the response writer is wrapped by a
// new *context which is passed to the next handlers ServeHTTP method.
func (c context) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&context{ResponseWriter: rw}, req)
	}
	return f
}
