package wraps

import (
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestGuardForbidden(t *testing.T) {
	h := wrap.New(
		Guard(String("forbidden")),
		wrap.Handler(String("hu?")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "forbidden", 200)
	if err != nil {
		t.Error(err)
	}
}

func TestGuardPassthrough(t *testing.T) {
	h := wrap.New(
		GuardFunc(wrap.NoOp),
		wrap.Handler(String("*")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "*", 200)
	if err != nil {
		t.Error(err)
	}
}
