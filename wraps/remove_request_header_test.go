package wraps

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func headerRequest() (wr *httptest.ResponseRecorder, rq *http.Request) {
	wr, rq = NewTestRequest("GET", "/")
	rq.Header.Set("X-test", "myval")
	return
}

func writeHeader(w http.ResponseWriter, rq *http.Request) {
	fmt.Fprint(w, rq.Header.Get("X-test"))
}

func TestRemoveRequestHeaderNoMatch(t *testing.T) {
	h := wrap.New(
		RemoveRequestHeader("Y-test"),
		wrap.HandlerFunc(writeHeader),
	)

	rw, rq := headerRequest()
	h.ServeHTTP(rw, rq)

	err := AssertResponse(rw, "myval", 200)
	if err != nil {
		t.Error(err)
	}

}

func TestRemoveRequestHeaderExactMatch(t *testing.T) {
	h := wrap.New(
		RemoveRequestHeader("X-test"),
		wrap.HandlerFunc(writeHeader),
	)

	rw, rq := headerRequest()
	h.ServeHTTP(rw, rq)
	err := AssertResponse(rw, "", 200)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveRequestHeaderPartialMatch(t *testing.T) {
	h := wrap.New(
		RemoveRequestHeader("X-"),
		wrap.HandlerFunc(writeHeader),
	)

	rw, rq := headerRequest()
	h.ServeHTTP(rw, rq)
	err := AssertResponse(rw, "", 200)
	if err != nil {
		t.Error(err)
	}
}
