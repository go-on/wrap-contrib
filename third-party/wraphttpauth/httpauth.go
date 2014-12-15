/*
Package wraphttpauth provides wrappers based on the github.com/abbot/go-http-auth package.

*/
package wraphttpauth

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/go-on/go-http-auth.v1"
	"gopkg.in/go-on/wrap.v2"
)

type digest struct {
	secrets func(user, realm string) string
	realm   string
}

// to be sure to implement the ContextWrapper interface
var _ wrap.ContextWrapper = &digest{}

// Validate makes sure that ctx supports the needed types
func (d *digest) ValidateContext(ctx wrap.Contexter) {
	var r auth.AuthenticatedRequest
	ctx.SetContext(&r)
	ctx.Context(&r)
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

// to be sure to implement the ContextWrapper interface
var _ wrap.ContextWrapper = &basic{}

// Validate makes sure that ctx supports the needed types
func (d *basic) ValidateContext(ctx wrap.Contexter) {
	var r auth.AuthenticatedRequest
	ctx.SetContext(&r)
	ctx.Context(&r)
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

func DigestSecret(user, password, realm string) string {
	m := md5.New()
	io.WriteString(m, user+":"+realm+":"+password)
	return fmt.Sprintf("%x", m.Sum(nil))
}

func BasicSecret(password, salt, magic string) string {
	return string(auth.MD5Crypt([]byte(password), []byte(salt), []byte(magic)))
	// md5e := auth.NewMD5Entry("$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1")
}
