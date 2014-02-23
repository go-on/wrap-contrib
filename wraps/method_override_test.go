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

func TestRest(t *testing.T) {
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
}
