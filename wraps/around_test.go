package wraps

import (
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestAround(t *testing.T) {
	h := wrap.New(
		Around(
			String("<body>"),
			String("</body>"),
		),
		AroundFunc(
			String("<h1>").ServeHTTP,
			String("</h1>").ServeHTTP,
		),
		wrap.Handler(String("rock around the clock")),
	)
	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "<body><h1>rock around the clock</h1></body>", 200)

	if err != nil {
		t.Error(err.Error())
	}
}
