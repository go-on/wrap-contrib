package main

import (
	"fmt"
	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"net/http"
)

type catcher struct{}

func (c catcher) Catch(p interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "catched: %s", p)
}

type panicker struct{}

func (p panicker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	panic("panic mysterious person found")
}

func main() {

	handler := wrap.New(
		wraps.HTMLContentType,
		wraps.GZip,
		wraps.ETag,
		wraps.Before(wraps.String(`<!DOCTYPE html><html lang="en"><head></head><body>`)),
		wraps.After(wraps.String(`</body></html>`)),
		wraps.Catch(catcher{}),
		wraps.Map(
			wraps.MatchQuery("name", "peter"), wraps.String("Hello Peter, how are you?"),
			wraps.MatchQuery("name", "mary"), wraps.String("Hello Mary, whats up?"),
			wraps.MatchQuery("name", "mister-x"), panicker{},
		),
		wraps.String(`
			<a href="/?name=peter">Peter</a><br />
			<a href="/?name=mary">Mary</a><br>
			<a href="/?name=mister-x">Mister X</a><br>
			`),
	)

	fmt.Println("go to localhost:8080")

	err := http.ListenAndServe(":8080", handler)

	if err != nil {
		println(err)
	}
}
