package wraps

import (
	"net/http"
	"testing"

	. "gopkg.in/go-on/wrap-contrib.v2/helper"

	"gopkg.in/go-on/wrap.v2"
)

func justSetHeader(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "text/plain")
}

func TestHeadRemoveBody(t *testing.T) {
	h := wrap.New(
		Head(),
		wrap.Handler(TextString("hi")),
	)

	rec, req := NewTestRequest("HEAD", "/")
	h.ServeHTTP(rec, req)
	err := AssertResponse(rec, "", 200)

	if err != nil {
		t.Error(err.Error())
	}

	ctype := rec.Header().Get("Content-Type")

	if ctype != "text/plain; charset=utf-8" {
		t.Errorf("Head should have Content-Type of text/plain; charset=utf-8, but has: %s", ctype)
	}
}

func TestHeadPassGet(t *testing.T) {
	h := wrap.New(
		Head(),
		wrap.Handler(TextString("hi")),
	)

	rec, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rec, req)
	err := AssertResponse(rec, "hi", 200)

	if err != nil {
		t.Error(err.Error())
	}

	ctype := rec.Header().Get("Content-Type")

	if ctype != "text/plain; charset=utf-8" {
		t.Errorf("Head should have Content-Type of text/plain; charset=utf-8, but has: %s", ctype)
	}
}

func TestHeadPassStatus(t *testing.T) {
	h := wrap.New(
		Head(),
		wrap.Handler(http.NotFoundHandler()),
	)

	rec, req := NewTestRequest("HEAD", "/")
	h.ServeHTTP(rec, req)
	err := AssertResponse(rec, "", 404)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestHeadNoChange(t *testing.T) {
	h := wrap.New(
		Before(http.HandlerFunc(justSetHeader)),
		Head(),
		wrap.HandlerFunc(wrap.NoOp),
	)

	rec, req := NewTestRequest("HEAD", "/")
	h.ServeHTTP(rec, req)
	err := AssertResponse(rec, "", 200)

	if err != nil {
		t.Error(err.Error())
	}

	ctype := rec.Header().Get("Content-Type")

	if ctype != "text/plain" {
		t.Errorf("Head should have Content-Type of text/plain, but has: %s", ctype)
	}
}
