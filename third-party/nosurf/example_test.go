package nosurf

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/go-on/wrap"
)

type context struct {
	http.ResponseWriter
	token Token
}

// context is an implementation for the Contexter interface.
//
// It receives a pointer to a value that is already stored inside the context.
// Values are distiguished by their type.
// Context sets the value of the given pointer to the value of the same type
// that is stored inside of the context.
// A pointer type that is not supported results in a panic.
func (c *context) Context(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *Token:
		*ty = c.token
	default:
		panic(fmt.Sprintf("unsupported context: %T", ctxPtr))
	}
}

// SetContext is an implementation for the Contexter interface.
//
// It receives a pointer to a value that will be stored inside the context.
// Values are distiguished by their type, that means that SetContext replaces
// and stored value of the same type.
// A pointer type that is not supported results in a panic.
func (c *context) SetContext(ctxPtr interface{}) {
	switch ty := ctxPtr.(type) {
	case *Token:
		c.token = *ty
	default:
		panic(fmt.Sprintf("unsupported context: %T", ctxPtr))
	}
}

// Wrap implements the wrap.Wrapper interface.
//
// When the request is served, the response writer is wrapped by a
// new *context which is passed to the next handlers ServeHTTP method.
func (c context) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&context{ResponseWriter: rw}, req)
	}
	return f
}

type app struct{}

// ServeHTTP serves the form value "a" for POST requests and otherwise the token
func (app) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		rw.Write([]byte(req.FormValue("a")))
		return
	}
	var token Token

	rw.(wrap.Contexter).Context(&token)
	rw.Write([]byte(string(token)))
	return
}

func Example() {
	stack := wrap.New(
		context{},
		&CheckToken{},
		SetToken{},
		wrap.Handler(app{}),
	)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	stack.ServeHTTP(rec, req)
	token := rec.Body.String()
	cookie := parseCookie(rec)

	rec = httptest.NewRecorder()
	req = mkPostReq(cookie, token)
	stack.ServeHTTP(rec, req)
	fmt.Println(rec.Code)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "http://localhost/", nil)
	stack.ServeHTTP(rec, req)
	token = rec.Body.String()
	cookie = parseCookie(rec)

	rec = httptest.NewRecorder()
	req = mkPostReq(cookie, token+"x")
	stack.ServeHTTP(rec, req)
	fmt.Println(rec.Code)
	fmt.Println(rec.Body.String())
	// Output:
	// 200
	// b
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
