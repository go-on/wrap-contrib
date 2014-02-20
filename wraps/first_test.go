package wraps

import (
	"net/http"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestFirstFirstWins(t *testing.T) {
	h := wrap.New(
		First(
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

func TestFirstSecondWins(t *testing.T) {
	h := wrap.New(
		First(
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

func TestFirstPassthrough(t *testing.T) {
	h := wrap.New(
		FirstFunc(
			DoNothing,
			DoNothing,
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
