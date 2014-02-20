package wraps

import (
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestGuardForbidden(t *testing.T) {
	h := wrap.New(
		Guard(Write("forbidden")),
		wrap.Handler(Write("hu?")),
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
		GuardFunc(DoNothing),
		wrap.Handler(Write("*")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "*", 200)
	if err != nil {
		t.Error(err)
	}
}
