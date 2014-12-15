package wraps

import (
	"net/http"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestFallbackFirstWins(t *testing.T) {
	h := wrap.New(
		Fallback(
			[]int{404},
			String("a"),
			String("b"),
		),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "a", 200)
	if err != nil {
		t.Error(err)
	}
}

func setHeader(w http.ResponseWriter, rq *http.Request) {
	// fmt.Printf("setting header for %s in %T\n", rq.URL.Path, w)
	w.Header().Set("x", "y")
	// if pk, ok := w.(*wrap.RWPeek); ok {
	// fmt.Printf("has changed %v\n", pk.HasChanged())
	// }
}

func TestFallbackFirstWinsSetHeaders(t *testing.T) {
	h := wrap.New(
		Fallback(
			[]int{404},
			http.HandlerFunc(setHeader),
			// String("a"),
			String("b"),
		),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	// fmt.Printf("headers %#v\n", rw.Header())
	if rw.Header().Get("x") != "y" {
		t.Errorf("header x should be y, but is %#v", rw.Header().Get("x"))
	}
}

func TestFallbackSecondWins(t *testing.T) {
	h := wrap.New(
		Fallback(
			[]int{404},
			http.Handler(http.NotFoundHandler()),
			String("b"),
		),
		wrap.Handler(String("*")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "b", 200)
	if err != nil {
		t.Error(err)
	}
}

func TestFallbackSecondWinsIgnoring(t *testing.T) {
	h := wrap.New(
		Fallback(
			[]int{404},
			http.Handler(http.NotFoundHandler()),
			String("b"),
		),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "b", 200)
	if err != nil {
		t.Error(err)
	}
}

func TestFallbackFirstWinsNotIgnoring(t *testing.T) {
	h := wrap.New(
		Fallback(
			[]int{405},
			http.Handler(http.NotFoundHandler()),
			String("b"),
		),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "404 page not found", 404)
	if err != nil {
		t.Error(err)
	}
}

func TestFallbackPassthrough(t *testing.T) {
	h := wrap.New(
		FallbackFunc(
			[]int{404},
			http.NotFoundHandler().ServeHTTP,
			http.NotFoundHandler().ServeHTTP,
		),
		wrap.Handler(String("*")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "*", 200)
	if err != nil {
		t.Error(err)
	}
}
