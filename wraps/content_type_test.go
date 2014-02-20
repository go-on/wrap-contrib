package wraps

import (
	"strings"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestContentTypeNothing(t *testing.T) {
	h := wrap.New(
		ContentType("my/contenttype"),
		wrap.HandlerFunc(DoNothing),
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
		Reader(strings.NewReader("hi")),
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
		wrap.HandlerFunc(NotFound),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)

	err := AssertResponse(rw, "not found", 404)
	if err != nil {
		t.Error(err.Error())
	}

	err = AssertHeader(rw, "Content-Type", "")
	if err != nil {
		t.Error(err.Error())
	}
}
