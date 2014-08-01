// Package nosurf provides integration of the github.com/justinas/nosurf package with the wrapper
// framework.

package nosurf

import (
	"net/http"

	"github.com/go-on/wrap"
	"github.com/justinas/nosurf"
)

type Token string

var TokenField = "csrf_token"

type SetToken struct{}

func (SetToken) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc

	f = func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			token := Token(nosurf.Token(req))
			rw.(wrap.Contexter).SetContext(&token)
		}
		next.ServeHTTP(rw, req)
	}
	return f
}

type CheckToken struct {
	FailureHandler http.Handler
	BaseCookie     *http.Cookie
	ExemptPaths    []string
	ExemptGlobs    []string
	ExemptRegexps  []interface{}
	ExemptFunc     func(r *http.Request) bool
}

func (c *CheckToken) Wrap(next http.Handler) http.Handler {
	ns := nosurf.New(next)
	if c.BaseCookie != nil {
		ns.SetBaseCookie(*c.BaseCookie)
	}
	if c.FailureHandler != nil {
		ns.SetFailureHandler(c.FailureHandler)
	}

	if len(c.ExemptPaths) > 0 {
		ns.ExemptPaths(c.ExemptPaths...)
	}

	if len(c.ExemptGlobs) > 0 {
		ns.ExemptPaths(c.ExemptGlobs...)
	}

	if len(c.ExemptRegexps) > 0 {
		ns.ExemptRegexps(c.ExemptRegexps...)
	}

	if c.ExemptFunc != nil {
		ns.ExemptFunc(c.ExemptFunc)
	}
	return ns
}
