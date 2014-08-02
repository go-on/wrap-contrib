// Package wraphttpauth provides integration of the github.com/abbot/go-http-auth package with the wrapper
// framework

package wraphttpauth

import (
	"net/http"

	"github.com/abbot/go-http-auth"
	"github.com/go-on/wrap"
)

/*
func Secret(user, realm string) string {
if user == "john" {
// password is "hello"
return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
}
return ""
}


func Secret(user, realm string) string {
if user == "john" {
// password is "hello"
return "b98e16cbc3d01734b264adba7baa3bf9"
}
return ""
}

	authenticator := auth.NewDigestAuthenticator("example.com", Secret)
http.HandleFunc("/", authenticator.Wrap(handle))


	authenticator := auth.NewBasicAuthenticator("example.com", Secret)
http.HandleFunc("/", authenticator.Wrap(handle))
*/

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

func Basic(realm string, secrets func(user, realm string) string) wrap.Wrapper {
	return &basic{secrets, realm}
}
