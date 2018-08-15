package wraps

import (
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestSetRequestHeader(t *testing.T) {
	h := wrap.New(
		SetRequestHeader("Y-test", "hiho"),
		wrap.Handler(String("huho")),
	)

	rw, rq := headerRequest()
	h.ServeHTTP(rw, rq)

	err := AssertResponse(rw, "huho", 200)
	if err != nil {
		t.Error(err)
	}

	header := rq.Header.Get("Y-test")
	if header != "hiho" {
		t.Errorf("header Y-test should be \"hiho\", but is %#v", header)
	}

}
