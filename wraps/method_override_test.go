package wraps

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
)

func methodWrite(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, req.Method)
}

func TestMethodOverride(t *testing.T) {
	r := wrap.New(
		MethodOverride(),
		wrap.HandlerFunc(methodWrite),
	)

	for overrideMethod, requestMethod := range acceptedOverrides {
		rw, req := helper.NewTestRequest(requestMethod, "/")
		req.Header.Set(overrideHeader, overrideMethod)
		r.ServeHTTP(rw, req)
		err := helper.AssertResponse(rw, overrideMethod, 200)
		if err != nil {
			t.Error(err.Error())
		}
		h := req.Header.Get(overrideHeader)
		if h != "" {
			t.Errorf("override header should be removed, but is %#v", h)
		}
	}

	rw, req := helper.NewTestRequest("GET", "/")
	req.Header.Set(overrideHeader, "GET")

	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw,
		"Unsupported value for X-HTTP-Method-Override: \"GET\".\nSupported values are PUT, DELETE, PATCH and OPTIONS",
		http.StatusPreconditionFailed)
	if err != nil {
		t.Error(err.Error())
	}

	rw, req = helper.NewTestRequest("GET", "/")
	req.Header.Set(overrideHeader, "PATCH")
	MethodOverride().ServeHTTP(rw, req)

	err = helper.AssertResponse(rw,
		"X-HTTP-Method-Override with value PATCH only allowed for POST requests.",
		http.StatusPreconditionFailed)
	if err != nil {
		t.Error(err.Error())
	}

}
