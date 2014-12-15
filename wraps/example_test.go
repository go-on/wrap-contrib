package wraps_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/wraps"
)

type catcher struct{}

func (c catcher) Catch(p interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "catched: %s", p)
}

type panicker struct{}

func (p panicker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	panic("panic mysterious person found")
}

func Example() {

	handler := wrap.New(

		wraps.HTMLContentType,            // sets the content type to text/html; charset=utf-8
		wraps.IfNoneMatch,                // sets 304 (not modified) for If-None-Match headers that matches the etag
		wraps.ETag,                       // calculates the ETag and sets the corresponding header
		wraps.Before(wraps.String(`<<`)), // writes << before everything further down the chain
		wraps.After(wraps.String(`>>`)),  // writes >> after everything further down the chain
		wraps.Catch(catcher{}),           // catches any panic and let catcher{} handle them

		// forwards the request to a set of request matcher and handler.
		// handler of the first matching matcher is used
		wraps.Map(
			wraps.MatchQuery("name", "peter"), wraps.String("Hi Peter!"),
			wraps.MatchQuery("name", "mary"), wraps.String("Hello Mary!"),
			wraps.MatchQuery("name", "mister-x"), panicker{},
		),
	)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?name=peter", nil)
	handler.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/?name=mary", nil)
	handler.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/?name=mister-x", nil)
	handler.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// Output:
	// <<Hi Peter!>>
	// <<Hello Mary!>>
	// <<catched: panic mysterious person found>>

}
