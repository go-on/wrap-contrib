// Package wraphttpauth provides integration of the github.com/abbot/go-http-auth package with the wrapper
// framework

package wraphttpauth

import (
	"net/http"

	"github.com/abbot/go-http-auth"
	"github.com/go-on/wrap"
)

type digest struct {
	secrets func(user, realm string) string
	realm   string
}

func (d *digest) Wrap(next http.Handler) http.Handler {
	authenticator := auth.NewDigestAuthenticator(d.realm, d.secrets)
	fn := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		w.(wrap.Contexter).SetContext(r)
		next.ServeHTTP(w, &r.Request)
	}
	return authenticator.Wrap(fn)
}

// Digest returns a wrapper that authenticates via auth.NewDigestAuthenticator
// and saves the resulting *auth.AuthenticatedRequest in the Contexter (response writer).
func Digest(realm string, secrets func(user, realm string) string) wrap.Wrapper {
	return &digest{secrets, realm}
}

type basic struct {
	secrets func(user, realm string) string
	realm   string
}

func (d *basic) Wrap(next http.Handler) http.Handler {
	authenticator := auth.NewBasicAuthenticator(d.realm, d.secrets)
	fn := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		w.(wrap.Contexter).SetContext(r)
		next.ServeHTTP(w, &r.Request)
	}
	return authenticator.Wrap(fn)
}

// Basic returns a wrapper that authenticates via auth.NewBasicAuthenticator
// and saves the resulting *auth.AuthenticatedRequest in the Contexter (response writer).
func Basic(realm string, secrets func(user, realm string) string) wrap.Wrapper {
	return &basic{secrets, realm}
}
