package wrapsession_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/stack"
	"gopkg.in/go-on/wrap-contrib.v2/third-party/wrapsession"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
	"gopkg.in/go-on/sessions.v1"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func errorHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "An error happened: %s\n", err.Error())
	fmt.Printf("An error happened: %s\n", err.Error())
	return
}

type writeToSession struct{}

var _ wrap.ContextWrapper = writeToSession{}

func (writeToSession) ValidateContext(ctx wrap.Contexter) {
	var session sessions.Session
	ctx.SetContext(&session)
	ctx.Context(&session)
}

func (writeToSession) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(w http.ResponseWriter, r *http.Request) {
		var session sessions.Session
		w.(wrap.Contexter).Context(&session)
		session.Values["name"] = r.URL.Query().Get("name")
		session.AddFlash("Hello, flash messages world!")
		next.ServeHTTP(w, r)
	}
	return f
}

type readFromSession struct{}

var _ wrap.ContextWrapper = readFromSession{}

func (readFromSession) ValidateContext(ctx wrap.Contexter) {
	var session sessions.Session
	ctx.SetContext(&session)
	ctx.Context(&session)
}

func (readFromSession) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(w http.ResponseWriter, r *http.Request) {
		var session sessions.Session
		w.(wrap.Contexter).Context(&session)
		fmt.Printf("Name: %s\n", session.Values["name"])
		for _, fl := range session.Flashes() {
			fmt.Printf("Flash: %v\n", fl)
		}
	}
	return f
}

func Example() {
	// stack.New() sets and checks global Contexter
	stack.New(&context{})

	// stack.Use() checks if global Contexter supports the context data needed by the wrappers
	// before adding them to the global Contexter
	stack.Use(
		wraps.ErrorHandlerFunc(errorHandler),
		wrapsession.SaveAndClear,
		wrapsession.Session(store, "my-session-name"),
		writeToSession{},
		readFromSession{},
	)

	// put the stack together and return a handler
	h := stack.Handler()

	req, _ := http.NewRequest("GET", "/?name=Peter", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	// Output:
	// Name: Peter
	// Flash: Hello, flash messages world!
}
