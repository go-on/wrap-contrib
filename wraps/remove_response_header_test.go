package wraps

import (
	"net/http"
	"testing"

	. "gopkg.in/go-on/wrap-contrib.v2/helper"

	"gopkg.in/go-on/wrap.v2"
)

func serveRemoveResponseTest(h http.Handler) (ctype string) {
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	ctype = rw.Header().Get("Content-Type")
	return
}

func TestRemoveResponseHeaderNoMatch1(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("X-"),
		wrap.Handler(TextString("hi")),
	)

	ctype := serveRemoveResponseTest(h)
	if ctype != "text/plain; charset=utf-8" {
		t.Errorf("wrong content type, should be text/plain; charset=utf-8, is: %s", ctype)
	}

}

func TestRemoveResponseHeaderNoMatch2(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("X-"),
		TextContentType,
	)

	ctype := serveRemoveResponseTest(h)
	if ctype != "text/plain; charset=utf-8" {
		t.Errorf("wrong content type, should be text/plain; charset=utf-8, is: %s", ctype)
	}

}

func TestRemoveResponseHeaderExactMatch1(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("Content-Type"),
		wrap.Handler(TextString("hi")),
	)
	ctype := serveRemoveResponseTest(h)
	if ctype != "" {
		t.Errorf("wrong content type, should be empty, is: %s", ctype)
	}
}

func TestRemoveResponseHeaderExactMatch2(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("Content-Type"),
		// wrap.Handler(String("hi")),
		HTMLContentType,
	)
	ctype := serveRemoveResponseTest(h)
	if ctype != "" {
		t.Errorf("wrong content type, should be empty, is: %s", ctype)
	}
}

func TestRemoveResponseHeaderPartialMatch(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("Content"),
		wrap.Handler(String("hi")),
	)

	ctype := serveRemoveResponseTest(h)
	if ctype != "" {
		t.Errorf("wrong content type, should be empty, is: %s", ctype)
	}
}
