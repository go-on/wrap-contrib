package wraps

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

type dispatchQ struct {
	method, path string
}

type ctx5 struct{ d, e string }

func (c *ctx5) black(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "black: d is %#v e is %#v", c.d, c.e)
}

func (c *ctx5) white(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "white: d is %#v e is %#v", c.d, c.e)
}

// note that this is not a pointer method, so every call is on a fresh instance
func (c ctx5) New(req *http.Request) http.HandlerFunc {
	q := req.URL.Query()
	c.d = q.Get("d")
	c.e = q.Get("e")

	switch req.URL.Path {
	case "/context/black":
		return c.black
	case "/context/white":
		return c.white
	default:
		return nil
	}
}

func TestDispatch(t *testing.T) {
	dispatchFnA := func(req *http.Request) http.Handler {
		if req.URL.Path == "/a" {
			return String("is a")
		}
		return nil
	}

	dispatchFnB := func(req *http.Request) http.Handler {
		if req.URL.Path == "/b" {
			return String("is b")
		}
		return nil
	}

	r := wrap.New(
		GETHandler("/company", String("get company")),
		POSTHandler("/company", String("post company")),
		PATCHHandler("/company", String("patch company")),
		DELETEHandler("/company", String("delete company")),
		HEADHandler("/company", String("head company")),
		PUTHandler("/company", String("put company")),
		OPTIONSHandler("/company", String("options company")),

		Map(
			// MatchScheme("http"), He

			MatchQuery("name", "peter"), String("peter"),
			MatchQuery("name", "paul"), String("paul"),
			And(MatchMethod("POST"), MatchPath("/hi")), String("ho"),
			MatchPath("/hi"), String("hi"),

			MatchPathRegex(regexp.MustCompile(`\/person\/customer\/[0-9]+`)), Map(
				MatchMethod("GET"),
				String("person customers"),
			),

			MatchPath("/blubb"),
			DispatchFunc(dispatchFnA),

			MatchPath("/person"), &MethodHandler{
				GET:     String("get person"),
				POST:    String("post person"),
				PATCH:   String("patch person"),
				DELETE:  String("delete person"),
				HEAD:    String("head person"),
				PUT:     String("put person"),
				OPTIONS: String("options person"),
			},
		),
		&MethodHandler{
			OPTIONS: String("my options"),
		},
		// StructDispatch(ctx5{}),
		DispatchFunc(dispatchFnA),
		Dispatch(DispatchFunc(dispatchFnB)),
		wrap.Handler(GETHandler("/hu", String("get hu"))),
		// wrap.Handler(String("not found")),
	)

	tests := map[dispatchQ]string{
		{"GET", "/?name=peter"}: "peter",
		{"GET", "/?name=paul"}:  "paul",
		{"POST", "/hi"}:         "ho",
		{"GET", "/hi"}:          "hi",

		{"GET", "/person"}:     "get person",
		{"POST", "/person"}:    "post person",
		{"PATCH", "/person"}:   "patch person",
		{"DELETE", "/person"}:  "delete person",
		{"HEAD", "/person"}:    "head person",
		{"PUT", "/person"}:     "put person",
		{"OPTIONS", "/person"}: "options person",
		{"TRACE", "/person"}:   "",

		{"GET", "/company"}:     "get company",
		{"POST", "/company"}:    "post company",
		{"PATCH", "/company"}:   "patch company",
		{"DELETE", "/company"}:  "delete company",
		{"HEAD", "/company"}:    "head company",
		{"PUT", "/company"}:     "put company",
		{"OPTIONS", "/company"}: "options company",

		{"GET", "/a"}:       "is a",
		{"GET", "/b"}:       "is b",
		{"GET", "/blubb"}:   "",
		{"OPTIONS", "/xyz"}: "my options",
		{"GET", "/hu"}:      "get hu",

		{"GET", "/xyz"}: "",

		{"GET", "/person/customer/6"}: "person customers",

		// dispatchQ{"GET", "/context/black?d=ddd&e=eee"}:   `black: d is "ddd" e is "eee"`,
		// dispatchQ{"GET", "/context/white?d=dddd&e=eeee"}: `white: d is "dddd" e is "eeee"`,

		// dispatchQ{"GET", "/context/black?d=ddd&e=eee"}:   ``,
		// dispatchQ{"GET", "/context/white?d=dddd&e=eeee"}: ``,
	}

	for q, res := range tests {
		// fmt.Printf("requesting %s %s\n", q.method, q.path)
		rw, req := NewTestRequest(q.method, q.path)
		r.ServeHTTP(rw, req)
		err := AssertResponse(rw, res, 200)
		if err != nil {
			t.Error(err)
		}
	}

}
