package stack

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

// basicContext is a simple basic context that supports the setting and getting of errors
type basicContext struct {
	http.ResponseWriter
	err error
}

// make sure to fulfill the ContextInjecter interface
var _ wrap.ContextInjecter = &basicContext{}
var _ = wrap.ValidateContextInjecter(&basicContext{})

func (c *basicContext) Context(ctxPtr interface{}) (found bool) {
	found = true
	switch ty := ctxPtr.(type) {
	case *error:
		if c.err == nil {
			return false
		}
		*ty = c.err
	case *http.ResponseWriter:
		*ty = c.ResponseWriter
	default:
		panic(&wrap.ErrUnsupportedContextGetter{ctxPtr})
	}
	return
}

func (c *basicContext) SetContext(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *error:
		c.err = *ty
	default:
		panic(&wrap.ErrUnsupportedContextSetter{ctxPtr})
	}
}

func (c basicContext) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&basicContext{ResponseWriter: rw}, req)
	}
	return f
}
