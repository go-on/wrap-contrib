package state

import (
	"fmt"
	"github.com/go-on/replacer"
	"gopkg.in/go-on/router.v2"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/account"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/basket"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/layout"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/login"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/loginstate"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/model"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/routes"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"net/http"
)

type State struct {
	user *model.User
}

func (s *State) forbidden(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("forbidden"))
}

func (s *State) checkLogin(r *http.Request) {
	cookie, err := r.Cookie("testlogin")

	if err != nil {
		return
	}

	email := cookie.Value
	u := model.FindUser(email)
	if u == nil {
		return
	}
	s.user = u
}

func unknown(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "route not yet defined: %s", req.URL.Path)
}

func (s *State) Dispatch(req *http.Request) http.Handler {
	currentPlaceholder := replacer.GetPlaceholder(req)
	currentRoute := router.GetRouteDefinition(req)

	switch currentPlaceholder {

	case layout.LoginState_.Name():
		ls := loginstate.LoginState{"", routes.LoginURL, routes.LogoutURL}

		if s.user != nil {
			ls.Name = s.user.Name
		}

		return ls

	case layout.Nav_.Name():
		return http.HandlerFunc(layout.NavHandler)

	case layout.Main_.Name():
		if currentRoute == routes.Login.OriginalPath {
			return http.HandlerFunc(login.FormHandler)
		}

		if s.user == nil {
			// s.user = model.FindUser("daisy@duck.com")
			return http.HandlerFunc(s.forbidden)
		}

		switch currentRoute {

		case routes.Account.OriginalPath:
			return http.HandlerFunc((&account.Account{s.user}).Show)

		case routes.Basket.OriginalPath:
			return &basket.Basket{s.user}

		default:
			return http.HandlerFunc(unknown)
		}

	default:
		return nil
	}

}

func (s State) New(req *http.Request) wraps.Dispatcher {
	ptr := &s
	ptr.checkLogin(req)
	if ptr.user == nil {
		ptr.user = model.FindUser("daisy@duck.com")
	}
	return ptr
}
