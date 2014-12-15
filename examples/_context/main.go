package main

import (
	"fmt"
	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"net/http"
)

type context struct{ sharedA, sharedB string }

func (c *context) handleA(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "handleA: sharedA is %#v sharedB is %#v pointer is %p", c.sharedA, c.sharedB, c)
}

func (c *context) handleB(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "handleB: sharedA is %#v sharedB is %#v pointer is %p", c.sharedA, c.sharedB, c)
}

func (c *context) notFound(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "not found")
}

// note that this is not a pointer method, so every call is on a fresh instance
func (c context) New(req *http.Request) http.HandlerFunc {
	c.sharedA = req.URL.Query().Get("a")
	c.sharedB = req.URL.Query().Get("b")

	switch req.URL.Path {
	case "/a":
		return c.handleA
	case "/b":
		return c.handleB
	default:
		return c.notFound
	}
}

func main() {
	http.ListenAndServe(
		":8283",
		wrap.New(
			wraps.StructDispatch(
				// note that this is not a pointer method, so every call on New() is on a fresh instance
				context{},
			),
		),
	)
}
