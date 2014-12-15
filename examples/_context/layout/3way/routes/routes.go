package routes

import (
	"gopkg.in/go-on/router.v2"
	"gopkg.in/go-on/router.v2/route"
)

var (
	Router      = router.New()
	Account     *route.Route
	Basket      *route.Route
	Login       *route.Route
	Logout      *route.Route
	LoginSubmit *route.Route

	AccountURL     string
	BasketURL      string
	LoginURL       string
	LogoutURL      string
	LoginSubmitURL string
)

func SetURLs() {
	AccountURL = router.MustURL(Account)
	BasketURL = router.MustURL(Basket)
	LoginURL = router.MustURL(Login)
	LogoutURL = router.MustURL(Logout)
	LoginSubmitURL = router.MustURL(LoginSubmit)
}
