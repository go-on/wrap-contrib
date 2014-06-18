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

type context struct{ d, e string }

func (c *context) black(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "black: d is %#v e is %#v", c.d, c.e)
}

func (c *context) white(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "white: d is %#v e is %#v", c.d, c.e)
}

// note that this is not a pointer method, so every call is on a fresh instance
func (c context) New(req *http.Request) http.HandlerFunc {
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
			return Write("is a")
		}
		return nil
	}

	dispatchFnB := func(req *http.Request) http.Handler {
		if req.URL.Path == "/b" {
			return Write("is b")
		}
		return nil
	}

	r := wrap.New(
		GETHandler("/company", Write("get company")),
		POSTHandler("/company", Write("post company")),
		PATCHHandler("/company", Write("patch company")),
		DELETEHandler("/company", Write("delete company")),
		HEADHandler("/company", Write("head company")),
		PUTHandler("/company", Write("put company")),
		OPTIONSHandler("/company", Write("options company")),

		Map(
			// MatchScheme("http"), He

			MatchQuery("name", "peter"), Write("peter"),
			MatchQuery("name", "paul"), Write("paul"),
			And(MatchMethod("POST"), MatchPath("/hi")), Write("ho"),
			MatchPath("/hi"), Write("hi"),

			MatchPathRegex(regexp.MustCompile(`\/person\/customer\/[0-9]+`)), Map(
				MatchMethod("GET"),
				Write("person customers"),
			),

			MatchPath("/blubb"),
			DispatchFunc(dispatchFnA),

			MatchPath("/person"), &MethodHandler{
				GET:     Write("get person"),
				POST:    Write("post person"),
				PATCH:   Write("patch person"),
				DELETE:  Write("delete person"),
				HEAD:    Write("head person"),
				PUT:     Write("put person"),
				OPTIONS: Write("options person"),
			},
		),
		&MethodHandler{
			OPTIONS: Write("my options"),
		},
		StructDispatch(context{}),
		DispatchFunc(dispatchFnA),
		Dispatch(DispatchFunc(dispatchFnB)),
		wrap.Handler(GETHandler("/hu", Write("get hu"))),
		// wrap.Handler(Write("not found")),
	)

	tests := map[dispatchQ]string{
		dispatchQ{"GET", "/?name=peter"}: "peter",
		dispatchQ{"GET", "/?name=paul"}:  "paul",
		dispatchQ{"POST", "/hi"}:         "ho",
		dispatchQ{"GET", "/hi"}:          "hi",

		dispatchQ{"GET", "/person"}:     "get person",
		dispatchQ{"POST", "/person"}:    "post person",
		dispatchQ{"PATCH", "/person"}:   "patch person",
		dispatchQ{"DELETE", "/person"}:  "delete person",
		dispatchQ{"HEAD", "/person"}:    "head person",
		dispatchQ{"PUT", "/person"}:     "put person",
		dispatchQ{"OPTIONS", "/person"}: "options person",
		dispatchQ{"TRACE", "/person"}:   "",

		dispatchQ{"GET", "/company"}:     "get company",
		dispatchQ{"POST", "/company"}:    "post company",
		dispatchQ{"PATCH", "/company"}:   "patch company",
		dispatchQ{"DELETE", "/company"}:  "delete company",
		dispatchQ{"HEAD", "/company"}:    "head company",
		dispatchQ{"PUT", "/company"}:     "put company",
		dispatchQ{"OPTIONS", "/company"}: "options company",

		dispatchQ{"GET", "/a"}:       "is a",
		dispatchQ{"GET", "/b"}:       "is b",
		dispatchQ{"GET", "/blubb"}:   "",
		dispatchQ{"OPTIONS", "/xyz"}: "my options",
		dispatchQ{"GET", "/hu"}:      "get hu",

		dispatchQ{"GET", "/xyz"}: "",

		dispatchQ{"GET", "/person/customer/6"}: "person customers",

		dispatchQ{"GET", "/context/black?d=ddd&e=eee"}:   `black: d is "ddd" e is "eee"`,
		dispatchQ{"GET", "/context/white?d=dddd&e=eeee"}: `white: d is "dddd" e is "eeee"`,
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
