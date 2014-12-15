package wraps

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	// "net/http"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

type ctx struct {
	http.ResponseWriter
	context string
}

func (c *ctx) SetContext(context interface{}) {
	c.context = context.(string)
}

func (c *ctx) Context(context interface{}) bool {
	//*context = *c.context
	ctx := context.(*string)
	*ctx = c.context
	return true
}

func WrapCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(&ctx{ResponseWriter: rw, context: "hello "}, req)
	})
}

func contextSetter(next http.Handler, rw http.ResponseWriter, req *http.Request) {
	var hello string
	rw.(wrap.Contexter).Context(&hello)
	rw.(wrap.Contexter).SetContext(hello + "world")
	next.ServeHTTP(rw, req)
}

func contextWriter(rw http.ResponseWriter, req *http.Request) {
	var hello string
	rw.(wrap.Contexter).Context(&hello)
	fmt.Fprint(rw, hello)
}

func TestGzipContext(t *testing.T) {
	handler := wrap.New(
		wrap.WrapperFunc(WrapCtx),
		GZip,
		HTMLContentType,
		wrap.NextHandlerFunc(contextSetter),
		wrap.HandlerFunc(contextWriter),
	)

	rw, req := NewTestRequest("GET", "/")
	req.Header.Set("Accept-Encoding", "gzip")
	handler.ServeHTTP(rw, req)

	reader, err := gzip.NewReader(rw.Body)
	if err != nil {
		t.Error(err.Error())
	}

	var b []byte
	b, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Error(err.Error())
	}

	if string(b) != "hello world" {
		t.Errorf("wrong body, expected: %#v, got: %#v", "hello world", string(b))
	}

	err = AssertHeader(rw, "Content-Encoding", "gzip")
	if err != nil {
		t.Error(err.Error())
	}

	err = AssertHeader(rw, "Content-Type", "text/html; charset=utf-8")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGzip(t *testing.T) {
	css := "body{ color: green;}"
	handler := wrap.New(
		GZip,
		CSSString(css),
	)

	rw, req := NewTestRequest("GET", "/")
	req.Header.Set("Accept-Encoding", "gzip")
	handler.ServeHTTP(rw, req)

	reader, err := gzip.NewReader(rw.Body)
	if err != nil {
		t.Error(err.Error())
	}

	var b []byte
	b, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Error(err.Error())
	}

	if string(b) != css {
		t.Errorf("wrong body, expected: %#v, got: %#v", css, string(b))
	}

	err = AssertHeader(rw, "Content-Encoding", "gzip")
	if err != nil {
		t.Error(err.Error())
	}

	err = AssertHeader(rw, "Content-Type", "text/css; charset=utf-8")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGzipSkipped(t *testing.T) {
	css := "body{ color: green;}"
	handler := wrap.New(
		GZip,
		CSSString(css),
	)

	rw, req := NewTestRequest("GET", "/")
	handler.ServeHTTP(rw, req)

	err := AssertResponse(rw, css, 200)
	if err != nil {
		t.Error(err.Error())
	}

	err = AssertHeader(rw, "Content-Encoding", "")
	if err != nil {
		t.Error(err.Error())
	}

	err = AssertHeader(rw, "Content-Type", "text/css; charset=utf-8")
	if err != nil {
		t.Error(err.Error())
	}
}
