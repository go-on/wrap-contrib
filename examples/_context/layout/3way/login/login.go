package login

import (
	. "gopkg.in/go-on/lib.v2/html/ht"
	"gopkg.in/go-on/router.v2"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/model"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/routes"
	"net/http"
	"time"
)

func FormHandler(w http.ResponseWriter, req *http.Request) {
	FormPost(router.MustURL(routes.LoginSubmit),
		DIV("EMail: ", InputText("email", Value_("donald@duck.com"))),
		DIV("Password: ", InputText("password", Value_("entenhausen"))),
		InputSubmit("login"),
	).WriteTo(w)
}

func LogoutHandler(w http.ResponseWriter, req *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    "testlogin",
		Value:   "",
		Path:    "/",
		Domain:  "localhost",
		Expires: expire,
		//expire.Format(time.UnixDate), 86400, true, true, "test=tcookie", []string{"test=tcookie"}
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, req, router.MustURL(routes.Account), 302)
}

func LoginSubmitHandler(w http.ResponseWriter, req *http.Request) {
	email := req.FormValue("email")
	password := req.FormValue("password")
	// fmt.Printf("email: %#v password: %#v\n", email, password)
	if email == "" || password == "" {
		return
	}
	u := model.FindUser(email)
	// u, has := users[email]
	if u == nil || u.Password != password {
		http.Redirect(w, req, router.MustURL(routes.Login), 302)
		return
	}
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    "testlogin",
		Value:   email,
		Path:    "/",
		Domain:  "localhost",
		Expires: expire,
		//expire.Format(time.UnixDate), 86400, true, true, "test=tcookie", []string{"test=tcookie"}
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, req, router.MustURL(routes.Account), 302)
}
