package wraps

import (
	"compress/gzip"
	"io/ioutil"
	// "net/http"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

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
