package wraps

import (
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
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
