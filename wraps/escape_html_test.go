package wraps

import (
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestEscapeHTML1(t *testing.T) {
	h := wrap.New(
		EscapeHTML,
		wrap.Handler(String(`abc<d>"e'f&g`)),
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

func TestEscapeHTML3(t *testing.T) {
	rw, req := NewTestRequest("GET", "/")
	EscapeHTML.WrapFunc(String(`abc<d>"e'f&g`).ServeHTTP).ServeHTTP(rw, req)

	err := AssertResponse(rw, `abc&lt;d&gt;&#34;e&#39;f&amp;g`, 200)

	if err != nil {
		t.Error(err.Error())
	}
}
