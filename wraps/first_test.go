package wraps

import (
	"net/http"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestFirstFirstWins(t *testing.T) {
	h := wrap.New(
		First(
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

func TestFirstSecondWins(t *testing.T) {
	h := wrap.New(
		First(
			http.HandlerFunc(wrap.NoOp),
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

func TestFirstPassthrough(t *testing.T) {
	h := wrap.New(
		FirstFunc(
			wrap.NoOp,
			wrap.NoOp,
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
