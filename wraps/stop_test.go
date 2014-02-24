package wraps

import (
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestStop(t *testing.T) {
	h := wrap.New(
		Stop(),
		wrap.Handler(Write("a")),
	)

	buf := NewResponseBuffer(nil)
	_, req := NewTestRequest("GET", "/")
	h.ServeHTTP(buf, req)

	if buf.Code != 0 {
		t.Errorf("status code should be 0, but is %d", buf.Code)
	}

	if buf.BodyString() != "" {
		t.Errorf("body should be empty, but is %#v", buf.BodyString())
	}
}
