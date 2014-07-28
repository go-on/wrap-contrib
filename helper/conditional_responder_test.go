package helper

/*
import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-on/wrap"
)

func writeHeader(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("a", "b")
}

func writeCode(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(407)
}

func writeCodeCreated(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(201)
}

func TestCheckResponseHeader(t *testing.T) {

	ck := NewCheckedResponseWriter(nil, nil)

	writeHeader(ck, nil)

	if !ck.HasChanged() {
		t.Errorf("should have changed, but has not")
	}

	if ck.Header().Get("a") != "b" {
		t.Errorf("header a should be b, but is: %#v", ck.Header().Get("a"))
	}
}

func TestCheckResponseCode(t *testing.T) {

	ck := NewCheckedResponseWriter(nil, nil)

	writeCode(ck, nil)

	if !ck.HasChanged() {
		t.Errorf("should have changed, but has not")
	}

	if ck.Code != 407 {
		t.Errorf("code should be 407, but is: %v", ck.Code)
	}
}

func TestCheckResponseIsOk1(t *testing.T) {
	ck := NewCheckedResponseWriter(nil, nil)
	DoNothing(ck, nil)

	if !ck.IsOk() {
		t.Errorf("should be ok when doing nothing, but is not")
	}
}

func TestCheckResponseIsOk2(t *testing.T) {
	ck := NewCheckedResponseWriter(nil, nil)
	writeCodeCreated(ck, nil)

	if !ck.IsOk() {
		t.Errorf("should be ok with code 201, but is not")
	}
}

func TestCheckResponseIsOk3(t *testing.T) {
	ck := NewCheckedResponseWriter(nil, nil)
	writeCode(ck, nil)

	if ck.IsOk() {
		t.Errorf("should not be ok with code 407, but is")
	}
}

func TestCheckWriteCodeTo(t *testing.T) {
	ckA := NewCheckedResponseWriter(nil, nil)
	ckB := NewCheckedResponseWriter(nil, nil)

	writeCode(ckA, nil)

	ckA.WriteCodeTo(ckB)

	if !ckB.HasChanged() {
		t.Errorf("should have changed, but has not")
	}

	if ckB.Code != 407 {
		t.Errorf("code should be 407, but is: %v", ckB.Code)
	}

	// don't write a second time
	ckB.Code = 0
	ckA.WriteCodeTo(ckB)

	if ckB.Code != 0 {
		t.Errorf("code should be 0, but is: %v", ckB.Code)
	}
}

func TestCheckWriteHeadersTo1(t *testing.T) {
	ckA := NewCheckedResponseWriter(nil, nil)
	ckB := NewCheckedResponseWriter(nil, nil)

	writeHeader(ckA, nil)

	ckA.WriteHeadersTo(ckB)

	if !ckB.HasChanged() {
		t.Errorf("should have changed, but has not")
	}

	if ckB.Header().Get("a") != "b" {
		t.Errorf("header a should be b, but is: %#v", ckB.Header().Get("a"))
	}

	// don't write a second time
	ckB.Header().Set("a", "")
	ckA.WriteHeadersTo(ckB)
	if ckB.Header().Get("a") != "" {
		t.Errorf(`header a should be "", but is: %#v`, ckB.Header().Get("a"))
	}
}

func TestCheckWriteHeadersTo2(t *testing.T) {
	ckA := NewCheckedResponseWriter(nil, nil)
	ckB := NewCheckedResponseWriter(nil, nil)

	writeHeader(ckA, nil)
	ckA.WriteCodeTo(ckB)

	defer func() {
		if recover() == nil {
			t.Errorf("should panic if code is written before headers, but does not")
		}
	}()

	ckA.WriteHeadersTo(ckB)

}

func TestCheckReset(t *testing.T) {

	ck := NewCheckedResponseWriter(nil, nil)

	writeHeader(ck, nil)
	writeCode(ck, nil)
	ck.Reset()

	if ck.HasChanged() {
		t.Errorf("should not have changed, but has")
	}

	if ck.Header().Get("a") != "" {
		t.Errorf(`header a should be "", but is: %#v`, ck.Header().Get("a"))
	}
}

func TestCheckWrite1(t *testing.T) {
	rec := httptest.NewRecorder()
	ck := NewCheckedResponseWriter(rec, nil)
	Write("hiho").ServeHTTP(ck, nil)

	if rec.Body.String() != "hiho" {
		t.Errorf(`body a should be "hiho", but is: %#v`, rec.Body.String())
	}
}

func TestCheckWrite2(t *testing.T) {
	rec := httptest.NewRecorder()
	ck := NewCheckedResponseWriter(rec, func(c *CheckedResponseWriter) bool {
		return true
	})
	Write("hiho").ServeHTTP(ck, nil)

	if rec.Body.String() != "hiho" {
		t.Errorf(`body a should be "hiho", but is: %#v`, rec.Body.String())
	}
}

func TestCheckWrite3(t *testing.T) {
	rec := httptest.NewRecorder()
	ck := NewCheckedResponseWriter(rec, func(c *CheckedResponseWriter) bool {
		return false
	})
	Write("hiho").ServeHTTP(ck, nil)

	if rec.Body.String() != "" {
		t.Errorf(`body a should be "", but is: %#v`, rec.Body.String())
	}
}

type ctx struct {
	http.ResponseWriter
	context string
}

func (c *ctx) SetContext(context interface{}) {
	c.context = context.(string)
}

func (c *ctx) Context(context interface{}) {
	//*context = *c.context
	ctx := context.(*string)
	*ctx = c.context
}

func contextSetter(rw http.ResponseWriter, req *http.Request) {
	var hello string
	rw.(wrap.ResponseWriterWithContext).Context(&hello)
	rw.(wrap.ResponseWriterWithContext).SetContext(hello + "world")
}

func TestCheckContext(t *testing.T) {
	c := &ctx{context: "hello "}
	ck := NewCheckedResponseWriter(c, nil)

	contextSetter(ck, nil)

	if c.context != "hello world" {
		t.Errorf(`context should be "hello world", but is: %#v`, c.context)
	}
}

func TestEscapeRWContext(t *testing.T) {
	c := &ctx{context: "hello "}
	esc := &EscapeHTMLResponseWriter{c}

	contextSetter(esc, nil)

	if c.context != "hello world" {
		t.Errorf(`context should be "hello world", but is: %#v`, c.context)
	}
}
*/
