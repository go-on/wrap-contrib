package wrapnosurf_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/third-party/wrapnosurf"
)

type app struct{}

var _ wrap.ContextWrapper = app{}

func (app) ValidateContext(ctx wrap.Contexter) {
	var token wrapnosurf.Token
	ctx.SetContext(&token)
	ctx.Context(&token)
}

// ServeHTTP serves the form value "a" for POST requests and otherwise the token
func (app) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			req.ParseForm()
			rw.Write([]byte(req.FormValue("a")))
			return
		}
		var token wrapnosurf.Token

		rw.(wrap.Contexter).Context(&token)
		rw.Write([]byte(string(token)))
		return
	}
	return f
}

func Example() {

	// ensure all context requirements are fulfilled
	wrap.ValidateWrapperContexts(&context{}, &wrapnosurf.CheckToken{}, wrapnosurf.SetToken{}, app{})

	stack := wrap.Stack(
		&context{},
		&wrapnosurf.CheckToken{},
		wrapnosurf.SetToken{},
		app{},
	)

	// here comes the tests
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	stack.ServeHTTP(rec, req)
	token := rec.Body.String()
	cookie := parseCookie(rec)

	rec = httptest.NewRecorder()
	req = mkPostReq(cookie, token)
	stack.ServeHTTP(rec, req)
	fmt.Println("-- success --")
	fmt.Println(rec.Code)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/", nil)
	stack.ServeHTTP(rec, req)
	fmt.Println("-- fail --")
	fmt.Println(rec.Code)
	fmt.Println(rec.Body.String())
	// Output:
	// -- success --
	// 200
	// b
	// -- fail --
	// 400
	//
}

func parseCookie(rec *httptest.ResponseRecorder) *http.Cookie {
	cookie := rec.Header().Get("Set-Cookie")
	cookie2 := cookie[0:strings.Index(cookie, ";")]
	splitter := strings.Index(cookie2, "=")
	c := http.Cookie{}
	c.Name = cookie2[0:splitter]
	c.Value = cookie2[splitter+1:]
	return &c
}

func mkPostReq(cookie *http.Cookie, token string) *http.Request {
	var vals url.Values = map[string][]string{}
	vals.Set("a", "b")
	req, _ := http.NewRequest("POST", "http://localhost/", strings.NewReader(vals.Encode()))
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("Referer", "http://localhost/")
	return req
}
