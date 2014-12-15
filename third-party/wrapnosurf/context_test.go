package wrapnosurf_test

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/third-party/wrapnosurf"
)

// context is an example how a wrap.Contexter can be build in order to get used
// with wrapnosurf
type context struct {
	http.ResponseWriter
	token wrapnosurf.Token
}

// make sure to fulfill the ContextInjecter interface
var _ wrap.ContextInjecter = &context{}
var _ = wrap.ValidateContextInjecter(&context{})

func (c *context) Context(ctxPtr interface{}) (found bool) {
	found = true
	switch ty := ctxPtr.(type) {
	case *wrapnosurf.Token:
		if string(c.token) == "" {
			return false
		}
		*ty = c.token
	case *http.ResponseWriter:
		*ty = c.ResponseWriter
	default:
		panic(&wrap.ErrUnsupportedContextGetter{ctxPtr})
	}
	return
}

func (c *context) SetContext(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *wrapnosurf.Token:
		c.token = *ty
	default:
		panic(&wrap.ErrUnsupportedContextSetter{ctxPtr})
	}
}

func (c context) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&context{ResponseWriter: rw}, req)
	}
	return f
}
