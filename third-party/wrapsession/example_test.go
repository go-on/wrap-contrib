package wrapsession_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-on/wrap-contrib/wraps"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/third-party/wrapsession"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func errorHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.(wrap.Contexter).Context(&err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "An error happened: %s\n", err.Error())
	fmt.Printf("An error happened: %s\n", err.Error())
}

func writeToSession(next http.Handler, w http.ResponseWriter, r *http.Request) {
	var session sessions.Session
	w.(wrap.Contexter).Context(&session)
	session.Values["name"] = r.URL.Query().Get("name")
	session.AddFlash("Hello, flash messages world!")
	next.ServeHTTP(w, r)
}

func readFromSession(w http.ResponseWriter, r *http.Request) {
	var session sessions.Session
	w.(wrap.Contexter).Context(&session)
	fmt.Printf("Name: %s\n", session.Values["name"])
	for _, fl := range session.Flashes() {
		fmt.Printf("Flash: %v\n", fl)
	}
}

func Example() {
	stack := wrap.New(
		context{},
		wraps.ErrorHandlerFunc(errorHandler),
		wrapsession.SaveAndClear,
		wrapsession.Session(store, "my-session-name"),
		wrap.NextHandlerFunc(writeToSession),
		wrap.HandlerFunc(readFromSession),
	)

	req, _ := http.NewRequest("GET", "/?name=Peter", nil)
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)

	// Output:
	// Name: Peter
	// Flash: Hello, flash messages world!
}
