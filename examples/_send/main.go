package main

import (
	"fmt"
	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"net/http"
	"net/http/httptest"
)

type path struct{}

func (path) Set(rw http.ResponseWriter, req *http.Request) (key string, val interface{}) {
	if req.URL.Path == "/" {
		return
	}
	return "path", req.URL.Path
}

type query struct{}

func (query) Set(rw http.ResponseWriter, req *http.Request) (key string, val interface{}) {
	return "query", req.URL.RawQuery
}

type r struct{}

func (r) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m, ok := wraps.Get(rw, "path").(string)
	if ok {
		fmt.Fprintf(rw, "got path: %#v\n", m)
	} else {
		fmt.Fprintf(rw, "got no path\n")
	}
	q, ok2 := wraps.Get(rw, "query").(string)
	if ok2 {
		fmt.Fprintf(rw, "got query: %#v\n", q)
	} else {
		fmt.Fprintf(rw, "got no query\n")
	}
}

func main() {
	h := wrap.New(
		wraps.Set(path{}),
		wraps.Set(query{}),
		wrap.Handler(r{}),
	)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/a?path=hu", nil)
	h.ServeHTTP(rec, req)
	println(rec.Body.String())
	req, _ = http.NewRequest("GET", "/?o=a", nil)
	rec.Body.Reset()
	h.ServeHTTP(rec, req)
	println(rec.Body.String())
}
