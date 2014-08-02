package wraps

import (
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
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
