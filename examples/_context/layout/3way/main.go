package main

import (
	"gopkg.in/go-on/router.v2"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/layout"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/login"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/routes"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/state"
	"net/http"
	"runtime"
	// _ "net/http/pprof"
	// "github.com/davecheney/profile"
	"os"
	"time"
)

func main() {
	// profile.CPUProfile
	// defer profile.Start(profile.CPUProfile).Stop()
	/*
		cfg := profile.Config{
			CPUProfile: true,
			//NoShutdownHook: true, // do not hook SIGINT
		}
	*/
	handler := layout.Layout.Template().DispatchFunc(state.State{}.New)

	routes.Account = routes.Router.GET("/", handler)
	routes.Basket = routes.Router.GET("/basket", handler)
	routes.Login = routes.Router.GET("/login", handler)
	routes.LoginSubmit = routes.Router.POSTFunc("/login-submit", login.LoginSubmitHandler)
	routes.Logout = routes.Router.GETFunc("/logout", login.LogoutHandler)

	router.Mount("/", routes.Router)
	routes.SetURLs()
	_ = time.Second
	server := &http.Server{
		Addr:    ":8283",
		Handler: nil,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 5 * time.Second,
	}
	//p := profile.Start(&cfg)
	_ = os.Exit

	go func() {
		time.Sleep(500 * time.Millisecond)
		runtime.GC()
	}()
	/*
		go func() {
			time.Sleep(15 * time.Second)
			println("I quit")
			p.Stop()
			os.Exit(0)
		}()
	*/
	server.ListenAndServe()
	// http.ListenAndServe(":8283", routes.Router)

}
