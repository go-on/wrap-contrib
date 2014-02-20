package wraps

import (
	"net/http"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

type stringTest struct{ text, contentType string }

func TestString(t *testing.T) {
	tests := map[stringTest]http.Handler{
		stringTest{"my text", "text/plain; charset=utf-8"}:           wrap.New(TextString("my text")),
		stringTest{"my css", "text/css; charset=utf-8"}:              wrap.New(CSSString("my css")),
		stringTest{"my json", "application/json; charset=utf-8"}:     wrap.New(JSONString("my json")),
		stringTest{"my js", "application/javascript; charset=utf-8"}: wrap.New(JavaScriptString("my js")),
		stringTest{"my html", "text/html; charset=utf-8"}:            wrap.New(HTMLString("my html")),
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
