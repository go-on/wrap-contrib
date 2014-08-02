package wraps

import (
	"net/http/httptest"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestStop(t *testing.T) {
	h := wrap.New(
		Stop(),
		wrap.Handler(String("a")),
	)

	rec := httptest.NewRecorder()
	_, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("status code should be 200, but is %d", rec.Code)
	}

	if rec.Body.String() != "" {
		t.Errorf("body should be empty, but is %#v", rec.Body.String())
	}
}
