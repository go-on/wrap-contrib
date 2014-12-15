package wraps

import (
	"net/http"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

type panicker struct{}

func (panicker) Catch(p interface{}, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(p.(string)))
}

func (panicker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("don't panic")
}

func TestCatchPanic(t *testing.T) {
	p := panicker{}
	h := wrap.New(
		Catch(p),
		wrap.Handler(p),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "don't panic", 200)
	if err != nil {
		t.Error(err)
	}
}

func TestCatchPanicCatchFunc(t *testing.T) {
	p := panicker{}
	h := wrap.New(
		// you should not do this and should simply use CatchFunc(p.Catch), its only for the test
		Catch(CatchFunc(p.Catch)),
		wrap.Handler(p),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "don't panic", 200)
	if err != nil {
		t.Error(err)
	}

}

func TestCatchNoPanic(t *testing.T) {
	p := panicker{}
	h := wrap.New(
		Catch(p),
		wrap.Handler(String("hi!")),
	)

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "hi!", 200)
	if err != nil {
		t.Error(err)
	}
}

func TestCatchNoPanicHeader(t *testing.T) {
	p := panicker{}
	h := wrap.New(
		Catch(p),
		wrap.HandlerFunc(setHeader),
		// wrap.Handler(Write("hi!")),
	)

	rw, req := NewTestRequest("GET", "/catch-no-panic-header")
	h.ServeHTTP(rw, req)

	if x := rw.Header().Get("x"); x != "y" {
		t.Errorf("header x should be y, but is %#v", x)
	}

}
