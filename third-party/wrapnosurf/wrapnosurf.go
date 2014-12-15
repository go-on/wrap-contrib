/*
Package wrapnosurf provides wrappers based on the github.com/justinas/nosurf package.

*/
package wrapnosurf

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/nosurf.v1"
)

// Token is the type that is saved inside a wrap.Contexter and
// represents a csrf token from the github.com/justinas/nosurf package.
type Token string

// Tokenfield is the name of the form field that submits a csrf token
var TokenField = "csrf_token"

// SetToken is a wrap.Wrapper that sets a csrf token in the Contexter
// (response writer) on GET requests.
type SetToken struct{}

var _ wrap.ContextWrapper = SetToken{}

func (SetToken) ValidateContext(ctx wrap.Contexter) {
	var token Token
	ctx.SetContext(&token)
	ctx.Context(&token)
}

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

// CheckToken is a wrap.Wrapper that checks the token via the github.com/justinas/nosurf
// package. Its attributes relate to the corresponding nosurf options. If they are nil,
// they are not set.
type CheckToken struct {
	FailureHandler http.Handler
	BaseCookie     *http.Cookie
	ExemptPaths    []string
	ExemptGlobs    []string
	ExemptRegexps  []interface{}
	ExemptFunc     func(r *http.Request) bool
}

// strictly not needed for the CheckToken wrapper, but since it makes no sense without
// SetToken, check it anyway
func (c *CheckToken) ValidateContext(ctx wrap.Contexter) {
	var token Token
	ctx.SetContext(&token)
	ctx.Context(&token)
}

var _ wrap.ContextWrapper = &CheckToken{}

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
