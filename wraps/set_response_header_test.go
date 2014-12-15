package wraps

import (
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestSetResponseHeader(t *testing.T) {
	h := wrap.New(
		SetResponseHeader("Y-test", "hiho"),
		wrap.Handler(String("huho")),
	)

	rw, rq := headerRequest()
	h.ServeHTTP(rw, rq)

	err := AssertResponse(rw, "huho", 200)
	if err != nil {
		t.Error(err)
	}

	err = AssertHeader(rw, "Y-test", "hiho")
	if err != nil {
		t.Error(err)
	}

}
