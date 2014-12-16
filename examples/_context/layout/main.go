package main

import (
	"fmt"
	. "gopkg.in/go-on/lib.v2/html/ht"
	"gopkg.in/go-on/lib.v2/html/types"
	"gopkg.in/go-on/lib.v2/html/types/placeholder"
	"github.com/go-on/replacer"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"net/http"
)

/*
external http.Handlers
*/
type titleHandler string

func (th titleHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "title: %s", th)
}

type bodyHandler string

func (bh bodyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "body: %s", bh)
}

/*
End of
external http.Handlers
*/

var title = placeholder.New(types.Text("title"))
var body = placeholder.New(types.Text("body"))

type page struct{ expensiveA, expensiveB string }

// this must be a pointer receiver to share the state
func (c *page) Dispatch(req *http.Request) http.Handler {
	switch replacer.GetPlaceholder(req) {
	case title.Name():
		s := fmt.Sprintf("expensiveA is %#v pointer is %p", c.expensiveA, c)
		return titleHandler(s)
	case body.Name():
		s := fmt.Sprintf("expensiveA is %#v expensiveB is %#v pointer is %p", c.expensiveA, c.expensiveB, c)
		return bodyHandler(s)
	default:
		return nil
	}
}

// note that this is not a pointer method, so every call is on a fresh instance
func (c page) New(req *http.Request) wraps.Dispatcher {
	c.expensiveA = req.URL.Query().Get("a")
	c.expensiveB = req.URL.Query().Get("b")
	return &c
}

func main() {
	handler := HTML5(
		TITLE(title),
		BODY(
			H1("Hello World"),
			P(body),
		),
	).Template().DispatchFunc(page{}.New)

	http.ListenAndServe(":8283", handler)
}
