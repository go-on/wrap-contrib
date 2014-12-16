package layout

import (
	. "gopkg.in/go-on/lib.v2/html/ht"
	"gopkg.in/go-on/lib.v2/html/types"
	"gopkg.in/go-on/lib.v2/html/types/placeholder"
	// "gopkg.in/go-on/router.v2"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/routes"
	"net/http"
)

var LoginState_ = placeholder.New(types.HTML("loginstate"))
var Main_ = placeholder.New(types.HTML("main"))
var Nav_ = placeholder.New(types.HTML("nav"))

var accountPath_ = placeholder.New(types.Attribute{"href", "accountPath"})
var basketPath_ = placeholder.New(types.Attribute{"href", "basketPath"})

var nav = UL(
	LI(A(accountPath_, "Account")),
	LI(A(basketPath_, "Basket")),
).Template("nav")

func NavHandler(rw http.ResponseWriter, req *http.Request) {
	nav.Replace(
		//accountPath_.Set(router.MustURL(routes.Account)),
		accountPath_.Set(routes.AccountURL),
		basketPath_.Set(routes.BasketURL),
		//basketPath_.Set(router.MustURL(routes.Basket)),
	).WriteTo(rw)
	/*
		UL(
			LI(AHref(router.MustURL(routes.Account), "Account")),
			LI(AHref(router.MustURL(routes.Basket), "Basket")),
		).WriteTo(rw)
	*/
}

var Layout = HTML5(
	TITLE("test"),
	BODY(
		LoginState_,
		BR(),
		Nav_,
		BR(),
		H1("Hello World"),
		DIV(Main_),
	),
)
