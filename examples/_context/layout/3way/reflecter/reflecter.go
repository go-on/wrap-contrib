package main

import (
	"fmt"
	"gopkg.in/go-on/lib.v2/html/element/compiler"
	. "gopkg.in/go-on/lib.v2/html/ht"
	"gopkg.in/go-on/router.v2"
	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"time"

	"net/http"
)

type page struct{ expensiveA, expensiveB string }

func (c *page) Title(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "title: expensiveA is %#v pointer is %p", c.expensiveA, c)
}

func (c *page) Body(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "body:  expensiveA is %#v expensiveB is %#v pointer is %p", c.expensiveA, c.expensiveB, c)
}

func (c *page) Prepare(rw http.ResponseWriter, req *http.Request) {
	c.expensiveA = req.URL.Query().Get("a")
	c.expensiveB = req.URL.Query().Get("b")
}

func currentTime(rw http.ResponseWriter, req *http.Request) {
	hour, min, sec := time.Now().Clock()
	fmt.Fprintf(rw, "<%02d:%02d:%02d>", hour, min, sec)
}

type body struct {
	path string
}

func (b *body) A(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "a: path is %#v pointer is %p", b.path, b)
}

func (b *body) B(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "b: path is %#v pointer is %p", b.path, b)
}

func (b *body) Prepare(w http.ResponseWriter, req *http.Request) {
	b.path = req.URL.Query().Get("path")
}

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "home")
}

func y(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "y")
}

func noLayout(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "no layout")
}

func main() {
	pg := compiler.NewDispatcher(&page{})
	bd := compiler.NewDispatcher(&body{})

	layoutRouter := router.New()
	layoutRouter.GETFunc("/", home)
	layoutRouter.GETFunc("/y", y)

	router.Mount("/", layoutRouter)

	nolayoutRouter := router.New()
	nolayoutRouter.GETFunc("/", noLayout)

	router.Mount("/nolayout", nolayoutRouter)

	layout := pg.DocHandler(
		//handler := compiler.DocHandler(
		HTML5(
			TITLE(pg.HTML("Title")),
			BODY(
				H1("Hello World"),
				P(pg.Text("Body")),
				P(wraps.EscapeHTML.Wrap(http.HandlerFunc(currentTime))),
				DIV(
					Style_("background-color:yellow;"),
					bd.ElementHandler("body",
						DIV(
							H3(bd.HTML("A")),
							H3(bd.Text("B")),
						),
					),
				),
				layoutRouter,
			),
		),
	)

	println("visit http://localhost:8283/?a=hu&b=ho&path=h<u>i</u>")

	http.ListenAndServe(":8283",
		wrap.New(
			wraps.Fallback([]int{405}, nolayoutRouter, layout),
		),
	)
}
