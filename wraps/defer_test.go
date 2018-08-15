package wraps

import (
	"net/http"
	"testing"

	"github.com/go-on/wrap"
	"gopkg.in/go-on/wrap-contrib.v2/helper"
)

func anyway(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`anyway`))
}

func TestDeferFunc(t *testing.T) {
	r := wrap.New(
		DeferFunc(anyway),
		wrap.Handler(panicker{}),
	)
	rw, req := helper.NewTestRequest("GET", "/")
	defer func() { recover() }()
	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw, "anyway", 200)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDefer(t *testing.T) {
	r := wrap.New(
		Defer(http.HandlerFunc(anyway)),
		wrap.Handler(panicker{}),
	)
	rw, req := helper.NewTestRequest("GET", "/")
	defer func() { recover() }()
	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw, "anyway", 200)
	if err != nil {
		t.Error(err.Error())
	}
}
