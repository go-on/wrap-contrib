package wraphttpauth_test

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"

	"gopkg.in/go-on/go-http-auth.v1"
)

// context is an example how a wrap.Contexter can be build in order to get used
// with wraphttpauth
type context struct {
	http.ResponseWriter
	authReq *auth.AuthenticatedRequest
}

// make sure to fulfill the ContextInjecter interface
var _ wrap.ContextInjecter = &context{}
var _ = wrap.ValidateContextInjecter(&context{})

func (c *context) Context(ctxPtr interface{}) (found bool) {
	found = true
	switch ty := ctxPtr.(type) {
	case *auth.AuthenticatedRequest:
		if c.authReq == nil {
			return false
		}
		*ty = *c.authReq
	case *http.ResponseWriter:
		*ty = c.ResponseWriter
	default:
		panic(&wrap.ErrUnsupportedContextGetter{ctxPtr})
	}
	return
}

func (c *context) SetContext(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *auth.AuthenticatedRequest:
		c.authReq = ty
	default:
		panic(&wrap.ErrUnsupportedContextSetter{ctxPtr})
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
