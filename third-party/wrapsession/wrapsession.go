/*
Package wrapsession provides wrappers based on the github.com/gorilla/sessions.
It assumes one session per request. If you need multiple sessions pre request, define your own
wrapper or use gorilla/sessions directly.

wrapsession expects the ResponseWriter to a wrap.Contexter supporting sessions.Session and error

*/
package wrapsession

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/context.v1"
	"gopkg.in/go-on/sessions.v1"
)

type session struct {
	Store sessions.Store
	Name  string
}

var _ wrap.ContextWrapper = &session{}

func (s *session) ValidateContext(ctx wrap.Contexter) {
	var sess sessions.Session
	ctx.SetContext(&sess)
	ctx.Context(&sess)
}

func (s *session) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		sess, _ := s.Store.Get(req, s.Name)
		if sess != nil {
			rw.(wrap.Contexter).SetContext(sess)
		}
		next.ServeHTTP(rw, req)
	}
	return f
}

func Session(store sessions.Store, name string) wrap.Wrapper {
	return &session{Store: store, Name: name}
}

type saveAndClear struct{}

var _ wrap.ContextWrapper = saveAndClear{}

func (s saveAndClear) ValidateContext(ctx wrap.Contexter) {
	var sess sessions.Session
	var err error
	ctx.SetContext(&sess)
	ctx.Context(&sess)
	ctx.SetContext(&err)
	ctx.Context(&err)
}

func (s saveAndClear) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			var sess sessions.Session
			ctx := rw.(wrap.Contexter)
			if ctx.Context(&sess) {
				err := sess.Save(req, rw)
				if err != nil {
					ctx.SetContext(&err)
				}
			}
		}()
		next.ServeHTTP(rw, req)
	}
	return context.ClearHandler(f)
}

// SaveAndClear saves the session and clears up any references to the request.
// it should be used at the beginning of the middleware chain (but after the Contexter)
var SaveAndClear = saveAndClear{}
