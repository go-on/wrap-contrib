package wraps

import (
	"errors"
	"net/http"
	"testing"

	. "gopkg.in/go-on/wrap-contrib.v2/helper"

	"gopkg.in/go-on/wrap.v2"
)

func TestErrorHandlerNoError1(t *testing.T) {
	stack := wrap.New(
		context{},
		ErrorHandler(String("error happend")),
		setError(false),
		String("all right"),
	)

	rec, req := NewTestRequest("GET", "/")
	stack.ServeHTTP(rec, req)

	err := AssertResponse(rec, "all right", 200)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestErrorHandlerNoError2(t *testing.T) {
	stack := wrap.New(
		context{},
		ErrorHandler(String("error happend")),
		setError(false),
		SetResponseHeader("a", "b"),
		// String("all right"),
	)

	rec, req := NewTestRequest("GET", "/")
	stack.ServeHTTP(rec, req)

	err := AssertHeader(rec, "a", "b")
	if err != nil {
		t.Error(err.Error())
	}
	err = AssertResponse(rec, "", 200)
	if err != nil {
		t.Error(err.Error())
	}
}

type setError bool

func (s setError) Wrap(next http.Handler) http.Handler {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		if s {
			err := errors.New("test error")
			rw.(wrap.Contexter).SetContext(&err)
		}
		next.ServeHTTP(rw, req)
	}
	return f
}

func TestErrorHandlerWithError(t *testing.T) {
	stack := wrap.New(
		context{},
		ErrorHandler(String("error happend")),
		setError(true),
		String("all right"),
	)

	rec, req := NewTestRequest("GET", "/")
	stack.ServeHTTP(rec, req)
	err := AssertResponse(rec, "error happend", 200)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestErrorHandlerFuncWithError(t *testing.T) {
	stack := wrap.New(
		context{},
		ErrorHandlerFunc(String("error happend").ServeHTTP),
		setError(true),
		String("all right"),
	)

	rec, req := NewTestRequest("GET", "/")
	stack.ServeHTTP(rec, req)
	err := AssertResponse(rec, "error happend", 200)
	if err != nil {
		t.Error(err.Error())
	}
}
