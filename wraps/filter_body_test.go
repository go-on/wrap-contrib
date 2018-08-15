package wraps

import (
	"net/http"
	"testing"

	"github.com/go-on/method"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
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
