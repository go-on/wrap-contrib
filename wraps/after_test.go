package wraps

import (
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestAfter(t *testing.T) {
	h := wrap.New(
		After(String(" after")),
		wrap.Handler(String("the day")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "the day after", 200)

	if err != nil {
		t.Error(err.Error())
	}
}
