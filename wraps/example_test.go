package wraps_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/wraps"
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
		wraps.HTMLContentType,
		wraps.ETag,
		wraps.Before(wraps.String(`<<`)),
		wraps.After(wraps.String(`>>`)),
		wraps.Catch(catcher{}),
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
