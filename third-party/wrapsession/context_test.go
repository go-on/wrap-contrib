package wrapsession_test

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/sessions.v1"
)

// context is an example how a wrap.Contexter can be build in order to get used
// with wraphttpauth
type context struct {
	http.ResponseWriter
	session *sessions.Session
	err     error
}

// make sure to fulfill the ContextInjecter interface
var _ wrap.ContextInjecter = &context{}
var _ = wrap.ValidateContextInjecter(&context{})

func (c *context) Context(ctxPtr interface{}) (found bool) {
	found = true
	switch ty := ctxPtr.(type) {
	case *sessions.Session:
		if c.session == nil {
			return false
		}
		*ty = *c.session
	case *http.ResponseWriter:
		*ty = c.ResponseWriter
	case *error:
		*ty = c.err
	default:
		panic(&wrap.ErrUnsupportedContextGetter{ctxPtr})
	}
	return
}

func (c *context) SetContext(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *sessions.Session:
		c.session = ty
	case *error:
		c.err = *ty
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
