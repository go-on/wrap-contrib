package wraps

import (
	"fmt"
	"strings"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	. "gopkg.in/go-on/wrap-contrib.v2/helper"
)

func TestReadSeeker(t *testing.T) {

	r := strings.NewReader("hiho")
	h := wrap.New(ReadSeeker(r))

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

func TestReadSeekerErrorRead(t *testing.T) {

	r := errorReader{}
	h := wrap.New(ReadSeeker(r))

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "", 200)

	if err == nil {
		t.Error("should get an error")
	}
}

func TestReadSeekerErrorSeek(t *testing.T) {
	r := errorReader{true}
	h := wrap.New(ReadSeeker(r))

	rw, req := NewTestRequest("GET", "/")
	h.ServeHTTP(rw, req)
	err := AssertResponse(rw, "", 200)

	if err == nil {
		t.Error("should get an error")
	}
}
