package wraps

import (
	"net/http"
	"testing"

	"gopkg.in/go-on/method.v1"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestFilterBody(t *testing.T) {
	r := wrap.New(
		FilterBody(method.GET),
		wrap.Handler(String("the body")),
	)

	rw, req := helper.NewTestRequest("GET", "/")

	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw,
		"",
		http.StatusOK)
	if err != nil {
		t.Error(err.Error())
	}

	rw, req = helper.NewTestRequest("POST", "/")
	r.ServeHTTP(rw, req)

	err = helper.AssertResponse(rw,
		"the body",
		http.StatusOK)
	if err != nil {
		t.Error(err.Error())
	}

	r = wrap.New(
		FilterBody(method.GET),
		wrap.Handler(http.NotFoundHandler()),
	)

	rw, req = helper.NewTestRequest("GET", "/")
	r.ServeHTTP(rw, req)
	err = helper.AssertResponse(rw,
		"",
		http.StatusNotFound)
	if err != nil {
		t.Error(err.Error())
	}

}
