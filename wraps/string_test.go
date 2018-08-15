package wraps

import (
	"net/http"
	"testing"

	"github.com/go-on/wrap"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

type stringTest struct{ text, contentType string }

func TestString(t *testing.T) {
	t.Skip()
	tests := map[stringTest]http.Handler{
		{"my text", "text/plain; charset=utf-8"}:           wrap.New(TextString("my text")),
		{"my css", "text/css; charset=utf-8"}:              wrap.New(CSSString("my css")),
		{"my json", "application/json; charset=utf-8"}:     wrap.New(JSONString("my json")),
		{"my js", "application/javascript; charset=utf-8"}: wrap.New(JavaScriptString("my js")),
		{"my html", "text/html; charset=utf-8"}:            wrap.New(HTMLString("my html")),
		{"just string", ""}:                                wrap.New(String("just string")),
	}

	for expected, handler := range tests {
		rw, req := NewTestRequest("GET", "/")
		handler.ServeHTTP(rw, req)
		err := AssertResponse(rw, expected.text, 200)
		if err != nil {
			t.Error(err.Error())
		}
		err = AssertHeader(rw, "Content-Type", expected.contentType)
		if err != nil {
			t.Error(err.Error())
		}
	}
}
