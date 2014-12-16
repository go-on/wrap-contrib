package account

import (
	. "gopkg.in/go-on/lib.v2/html/ht"
	"gopkg.in/go-on/lib.v2/html/types"
	"gopkg.in/go-on/lib.v2/html/types/placeholder"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/model"
	"net/http"
)

type Account struct{ User *model.User }

var userName = placeholder.New(types.HTML("username"))
var userEMail = placeholder.New(types.HTML("useremail"))

var show = TABLE(
	TR(TH("Name"), TH("EMail")),
	TR(TD(userName), TD(userEMail)),
).Template("account")

func (a *Account) Show(w http.ResponseWriter, req *http.Request) {
	show.Replace(
		userEMail.Set(a.User.EMail),
		userName.Set(a.User.Name),
	).WriteTo(w)
	/*
		TABLE(
			TR(TH("Name"), TH("EMail")),
			TR(TD(a.User.Name), TD(a.User.EMail)),
		).WriteTo(w)
	*/
}
