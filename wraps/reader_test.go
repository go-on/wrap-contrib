package wraps

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-on/wrap"
	. "github.com/go-on/wrap-contrib/helper"
)

func TestReader(t *testing.T) {

	r := strings.NewReader("hiho")
	h := wrap.New(Reader(r))

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "hiho", 200)

	if err != nil {
		t.Error(err.Error())
	}

	// read it a second time
	rw, req = NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err = AssertResponse(rw, "hiho", 200)

	if err != nil {
		t.Error(err.Error())
	}
}

type errorReader struct {
	errorOnSeek bool
}

func (errorReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("can't read")
}

func (e errorReader) Seek(int64, int) (i int64, err error) {
	if e.errorOnSeek {
		return 0, fmt.Errorf("can't seek")
	}
	return
}

func TestReaderErrorRead(t *testing.T) {

	r := errorReader{}
	h := wrap.New(Reader(r))

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "", 200)

	if err == nil {
		t.Error("should get an error")
	}
}

func TestReaderErrorSeek(t *testing.T) {
	r := errorReader{true}
	h := wrap.New(Reader(r))

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "", 200)

	if err == nil {
		t.Error("should get an error")
	}
}
