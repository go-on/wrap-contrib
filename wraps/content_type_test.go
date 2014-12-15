package wraps

import (
	"strings"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestContentTypeNothing(t *testing.T) {
	h := wrap.New(
		ContentType("my/contenttype"),
		wrap.Handler(TextString("hu")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertHeader(rw, "Content-Type", "my/contenttype")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestContentTypeOk(t *testing.T) {
	h := wrap.New(
		ContentType("my/contenttype"),
		ReadSeeker(strings.NewReader("hi")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)

	err := AssertResponse(rw, "hi", 200)
	if err != nil {
		t.Error(err.Error())
	}

	err = AssertHeader(rw, "Content-Type", "my/contenttype")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestContentTypeError(t *testing.T) {
	h := wrap.New(
		ContentType("my/contenttype"),
		wrap.HandlerFunc(WriteError),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)

	err := AssertResponse(rw, "500 server error", 500)
	if err != nil {
		t.Error(err.Error())
	}

	err = AssertHeader(rw, "Content-Type", "")
	if err != nil {
		t.Error(err.Error())
	}
}
