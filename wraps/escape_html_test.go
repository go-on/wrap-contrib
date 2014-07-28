package wraps

import (
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestEscapeHTML1(t *testing.T) {
	h := wrap.New(
		EscapeHTML,
		wrap.Handler(Write(`abc<d>"e'f&g`)),
	)

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, `abc&lt;d&gt;&#34;e&#39;f&amp;g`, 200)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestEscapeHTML2(t *testing.T) {
	h := wrap.New(
		EscapeHTML,
		HTMLContentType,
	)

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)

	expected := "text/html; charset=utf-8"
	got := rw.Header().Get("Content-Type")

	if got != expected {
		t.Errorf("expected: %#v, got: %#v", expected, got)
	}
}
