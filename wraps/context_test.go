package wraps

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

type context struct {
	http.ResponseWriter
	err error
}

// make sure to fulfill the ContextInjecter interface
var _ wrap.ContextInjecter = &context{}
var _ = wrap.ValidateContextInjecter(&context{})

func (c *context) Context(ctxPtr interface{}) (found bool) {
	found = true // save work
	switch ty := ctxPtr.(type) {
	case *http.ResponseWriter:
		*ty = c.ResponseWriter
	case *error:
		if c.err == nil {
			return false
		}
		*ty = c.err
	default:
		panic(&wrap.ErrUnsupportedContextGetter{ctxPtr})
	}
	return
}

func (c *context) SetContext(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *error:
		c.err = *ty
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
