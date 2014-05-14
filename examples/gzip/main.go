package main

import (
	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/wraps"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8081", wrap.New(
		wraps.GZip,
		wraps.CSSString("body { color: green; }"),
	))

	if err != nil {
		panic(err.Error())
	}
}
