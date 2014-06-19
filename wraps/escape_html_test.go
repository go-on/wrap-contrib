package wraps

import (
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestEscapeHTML(t *testing.T) {
	h := wrap.New(
		EscapeHTML,
		wrap.Handler(Write(`abc<d>"e'f&g`)),
	)

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, `abc&lt;d&gt;&#34;e&#39;f&amp;g`, 200)

	if err != nil {
		t.Error(err.Error())
	}
}
