package main

import (
	"fmt"

	"gopkg.in/go-on/router.v2"
	"gopkg.in/go-on/sessions.v1"
	// "gopkg.in/go-on/wrap.v2-contrib-testing/wrapstesting"
	"gopkg.in/go-on/context.v1"
	// "gopkg.in/go-on/router.v2"
	// "gopkg.in/go-on/router.v2/route"
	"net/http"
	"net/http/httptest"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
)

type contextkey int

const ctxUser contextkey = 0
const ctxSessionStore contextkey = 1
const ctxSessionName contextkey = 2

func GetSessionStore(r *http.Request) sessions.Store {
	if rv := context.Get(r, ctxSessionStore); rv != nil {
		return rv.(sessions.Store)
	}
	return nil
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	store := GetSessionStore(r)
	if store == nil {
		return nil, fmt.Errorf("no session store found")
	}
	name := context.Get(r, ctxSessionName)
	if name == nil {
		return nil, fmt.Errorf("no session name found")
	}
	return store.Get(r, name.(string))
}

func MustGetSession(r *http.Request) *sessions.Session {
	s, err := GetSession(r)
	if err != nil {
		panic(err.Error())
	}
	return s
}

func SetSessionStore(r *http.Request, s sessions.Store) {
	context.Set(r, ctxSessionStore, s)
}

func SetSessionName(r *http.Request, name string) {
	context.Set(r, ctxSessionName, name)
}

func GetUser(r *http.Request) *user {
	if rv := context.Get(r, ctxUser); rv != nil {
		return rv.(*user)
	}
	return nil
}

func SetUser(r *http.Request, u *user) {
	context.Set(r, ctxUser, u)
}

type delegate struct {
	a, b  string
	store sessions.Store
}

func (c *delegate) A(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "A: a is %#v", c.a)
}

func (c *delegate) B(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "B:  a is %#v b is %#v", c.a, c.b)
}

func (c delegate) Dispatch(req *http.Request) http.Handler {
	c.a = req.URL.Query().Get("a")
	c.b = req.URL.Query().Get("b")
	switch req.URL.Query().Get("target") {
	case "a":
		return http.HandlerFunc((&c).A)
	case "b":
		return http.HandlerFunc((&c).B)
	}
	return nil //http.HandlerFunc((&c).B)
}

func a(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<a>")
}

func b(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "<b>"+MustGetSession(req).Values["hi"].(string))
}

type setStore struct {
	store       sessions.Store
	sessionName string
}

func (sb *setStore) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	SetSessionStore(req, sb.store)
	SetSessionName(req, sb.sessionName)
}

type user struct {
	FirstName, LastName string
}

func s(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "<b>")
	SetUser(req, &user{"Donald", "Duck"})
	MustGetSession(req).Values["hi"] = "ho"
	// GetSession(req, "sessionName")
	// wraps.SetContext(req, "user", &user{"Donald", "Duck"})
}

func printUser(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "<b>")
	println("print user called")
	// var u user
	//	wraps.GetContext(req, "user", &u)
	u := GetUser(req)
	fmt.Fprintf(w, "Firstname: %s - Lastname: %s", u.FirstName, u.LastName)
	//wraps.SetContextJSON(req, "user", &user{"Donald", "Duck"})
}

// sessions.Save is no http.HandlerFunc (reversed parameters), we need this func
func saveSession(w http.ResponseWriter, req *http.Request) {
	sessions.Save(req, w)
}

//var store sessions.Store = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
	rt := router.New()
	rt.GETFunc("/a", a)
	rt.GETFunc("/b", b)

	rt.Mount("/", http.DefaultServeMux)
	store := sessions.NewCookieStore([]byte("something-very-secret"))

	h := wrap.New(
		// always clear the context
		wrap.WrapperFunc(context.ClearHandler),
		// always save the sessions
		wraps.DeferFunc(saveSession),
		wraps.Before(&setStore{store, "user"}),
		wraps.BeforeFunc(s),
		wraps.Around(wraps.String("<p>"), wraps.String("</p>")),
		wraps.EscapeHTML,
		wraps.Dispatch(rt),
		//wraps.Dispatch(delegate{}),
		wraps.Dispatch(&delegate{store: store}),
		wrap.HandlerFunc(printUser),
	)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/a?a=AAA&b=BBBB&target=a", nil)
	h.ServeHTTP(rec, req)
	println(rec.Body.String())
	req, _ = http.NewRequest("GET", "/?a=A&b=B&target=b", nil)
	rec.Body.Reset()
	h.ServeHTTP(rec, req)
	println(rec.Body.String())
	req, _ = http.NewRequest("GET", "/c?a=A&b=B", nil)
	rec.Body.Reset()
	h.ServeHTTP(rec, req)
	println(rec.Body.String())
}
