package wraps

import (
	"net/http"
	"testing"

	. "github.com/go-on/wrap-contrib/helper"

	"github.com/go-on/wrap"
)

func serveRemoveResponseTest(h http.Handler) (ctype string) {
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	ctype = rw.Header().Get("Content-Type")
	return
}

func TestRemoveResponseHeaderNoMatch(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("X-"),
		wrap.Handler(Write("hi")),
	)

	ctype := serveRemoveResponseTest(h)
	if ctype != "text/plain" {
		t.Errorf("wrong content type, should be text/plain, is: %s", ctype)
	}

}

func TestRemoveResponseHeaderExactMatch(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("Content-Type"),
		wrap.Handler(Write("hi")),
	)
	ctype := serveRemoveResponseTest(h)
	if ctype != "" {
		t.Errorf("wrong content type, should be empty, is: %s", ctype)
	}
}

func TestRemoveResponseHeaderPartialMatch(t *testing.T) {
	h := wrap.New(
		RemoveResponseHeader("Content"),
		wrap.Handler(Write("hi")),
	)

	ctype := serveRemoveResponseTest(h)
	if ctype != "" {
		t.Errorf("wrong content type, should be empty, is: %s", ctype)
	}
}
