package wraps

import (
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestBefore(t *testing.T) {
	h := wrap.New(
		Before(String("before ")),
		wrap.Handler(String("midnight")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "before midnight", 200)

	if err != nil {
		t.Error(err.Error())
	}
}
