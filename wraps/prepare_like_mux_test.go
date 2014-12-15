package wraps

import (
	"net/http"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestPrepareLikeMuxCleanPath(t *testing.T) {
	h := wrap.New(
		PrepareLikeMux(),
		wrap.Handler(String("one two")),
	)
	rw, req := NewTestRequest("GET", "/hi//")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "<a href=\"/hi/\">Moved Permanently</a>.", 301)

	if err != nil {
		t.Error(err.Error())
	}

	rw, req = NewTestRequest("GET", "")
	h.ServeHTTP(rw, req)
	err = AssertResponse(rw, "<a href=\"/\">Moved Permanently</a>.", 301)

	if err != nil {
		t.Error(err.Error())
	}

	rw, req = NewTestRequest("GET", "hi")
	h.ServeHTTP(rw, req)
	err = AssertResponse(rw, "<a href=\"/hi\">Moved Permanently</a>.", 301)

	if err != nil {
		t.Error(err.Error())
	}

	rw, req = NewTestRequest("GET", "/hi")
	h.ServeHTTP(rw, req)
	err = AssertResponse(rw, "one two", 200)

	if err != nil {
		t.Error(err.Error())
	}
}

type setRequestURI string

func (s setRequestURI) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.RequestURI = string(s)
		next.ServeHTTP(rw, req)
	})
}

func TestPrepareLikeMux(t *testing.T) {
	h := wrap.New(
		setRequestURI("*"),
		PrepareLikeMux(),
		wrap.Handler(String("one two")),
	)
	rw, req := NewTestRequest("GET", "/")
	req.ProtoMinor = 0

	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "", http.StatusBadRequest)

	if err != nil {
		t.Error(err.Error())
	}

	rw, req = NewTestRequest("GET", "/")

	h.ServeHTTP(rw, req)
	if rw.Header().Get("Connection") != "close" {
		t.Errorf("connection should be close but is: %#v", rw.Header().Get("Connection"))
	}

}
