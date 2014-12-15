package main

import (
	"fmt"
	// auth "github.com/go-on/go-http-auth"
	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"net/http"
)

/*
Authenticator interface {
	Password(user, realm string) (password string)
	Salt() string
	Magic() string
}
*/

type authenicator struct{}

func (a authenicator) Password(user, realm string) (password string) {
	var users = map[string]string{
		"a":    "b",
		"wiki": "pedia",
	}

	fmt.Printf("user: %#v, password: %#v\n", user, users[user])
	//s := string(auth.MD5Crypt([]byte(users[user]), []byte("hiho"), []byte("$huho$")))
	//fmt.Println(s)
	return users[user]
}

func (a authenicator) Salt() string  { return "hiho" }
func (a authenicator) Magic() string { return "huho" }

func form(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(`Loginform`))
}

func main() {
	f := http.HandlerFunc(form)
	_ = f
	h := wrap.New(
		wraps.BasicAuthentication("admin", authenicator{}, nil),
		wraps.HTMLString("<html><body>Welcome</body></html>"),
	)

	http.ListenAndServe(":8182", h)
}
