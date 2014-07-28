package wraps

import (
	"net/http"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestFallbackFirstWins(t *testing.T) {
	h := wrap.New(
		Fallback(
			[]int{404},
			Write("a"),
			Write("b"),
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
			// Write("a"),
			Write("b"),
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
			http.HandlerFunc(DoNothing),
			Write("b"),
		),
		wrap.Handler(Write("*")),
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
			http.HandlerFunc(NotFound),
			Write("b"),
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
			http.HandlerFunc(NotFound),
			Write("b"),
		),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "not found", 404)
	if err != nil {
		t.Error(err)
	}
}

func TestFallbackPassthrough(t *testing.T) {
	h := wrap.New(
		FallbackFunc(
			[]int{404},
			NotFound,
			NotFound,
		),
		wrap.Handler(Write("*")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "*", 200)
	if err != nil {
		t.Error(err)
	}
}
