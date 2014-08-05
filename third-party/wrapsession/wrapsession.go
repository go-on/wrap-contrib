/*
Package wrapsession provides wrappers based on the github.com/gorilla/sessions.
It assumes one session per request. If you need multiple sessions pre request, define your own
wrapper or use gorilla/sessions directly.

wrapsession expects the ResponseWriter to a wrap.Contexter supporting sessions.Session and error

*/
package wrapsession

import (
	"net/http"

	"github.com/go-on/wrap"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

type session struct {
	Store sessions.Store
	Name  string
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