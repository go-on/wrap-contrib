package loginstate

import (
	. "gopkg.in/go-on/lib.v3/html/ht"
	"gopkg.in/go-on/lib.v3/html/types"
	"gopkg.in/go-on/lib.v3/html/types/placeholder"
	"net/http"
)

type LoginState struct {
	Name       string
	LoginLink  string
	LogoutLink string
}

var loginLink_ = placeholder.New(types.Attribute{"href", "loginLink"})
var logoutLink_ = placeholder.New(types.Attribute{"href", "logoutLink"})
var name_ = placeholder.New(types.HTML("name"))

var loginStateLoggedOut = DIV(types.Class("login-state"), A(loginLink_, "login")).Template("logged-out")
var loginStateLoggedIn = DIV(types.Class("login-state"), "logged in as ", name_, " ", A(logoutLink_, "logout")).Template("logged-in")

func (l LoginState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// d := DIV(types.Class("login-state"))
	if l.Name == "" {
		loginStateLoggedOut.Replace(
			loginLink_.Set(l.LoginLink),
		).WriteTo(w)
		return
		// d.Add(AHref(l.LoginLink, "login"))
	}
	/*
		d.Add(
			"logged in as "+l.Name+" ",
			AHref(l.LogoutLink, "logout"),
		)
	*/
	loginStateLoggedIn.Replace(
		logoutLink_.Set(l.LogoutLink),
		name_.Set(l.Name),
	).WriteTo(w)
	//d.WriteTo(w)
}
