package helper

import (
	"net/http"
	"testing"
)

func TestResponseBufferWriteTo(t *testing.T) {
	buf := NewResponseBuffer(nil)
	rec, req := NewTestRequest("GET", "/")
	Write("hi").ServeHTTP(buf, req)
	buf.WriteTo(rec)
	err := AssertResponse(rec, "hi", 200)
	if err != nil {
		t.Error(err)
	}
}

func TestResponseBufferWriteToStatus(t *testing.T) {
	buf := NewResponseBuffer(nil)
	rec, req := NewTestRequest("GET", "/")
	NotFound(buf, req)
	buf.WriteTo(rec)
	err := AssertResponse(rec, "not found", 404)
	if err != nil {
		t.Error(err)
	}

	if buf.IsOk() {
		t.Error("buf is ok, but should be not")
	}
}

func TestResponseBufferChanged(t *testing.T) {
	buf1 := NewResponseBuffer(nil)
	buf2 := NewResponseBuffer(nil)
	_, req := NewTestRequest("GET", "/")
	Write("hi").ServeHTTP(buf1, req)
	buf1.WriteTo(buf2)

	if buf1.BodyString() != "hi" {
		t.Errorf("body string of buf1 should be \"hi\" but is :%#v", buf1.BodyString())
	}

	if buf2.BodyString() != "hi" {
		t.Errorf("body string of buf2 should be \"hi\" but is :%#v", buf2.BodyString())
	}

	if string(buf1.Body()) != "hi" {
		t.Errorf("body of buf1 should be \"hi\" but is :%#v", string(buf1.Body()))
	}

	if string(buf2.Body()) != "hi" {
		t.Errorf("body of buf2 should be \"hi\" but is :%#v", string(buf2.Body()))
	}

	if buf1.Code != 0 {
		t.Errorf("Code of buf1 should be %d but is :%d", 0, buf1.Code)
	}

	if buf2.Code != 0 {
		t.Errorf("Code of buf2 should be %d but is :%d", 0, buf2.Code)
	}

	ctype1 := buf1.Header().Get("Content-Type")
	if ctype1 != "text/plain" {
		t.Errorf("Content-Type of buf1 should be %#v but is: %#v", "text/plain", ctype1)
	}

	ctype2 := buf2.Header().Get("Content-Type")
	if ctype2 != "text/plain" {
		t.Errorf("Content-Type of buf2 should be %#v but is: %#v", "text/plain", ctype2)
	}

	if !buf1.HasChanged() {
		t.Error("buf1 should be changed, but is not")
	}

	if !buf2.HasChanged() {
		t.Error("buf2 should be changed, but is not")
	}

	if !buf1.IsOk() {
		t.Error("buf1 should be ok, but is not")
	}

	if !buf2.IsOk() {
		t.Error("buf2 should be ok, but is not")
	}
}

func TestResponseBufferNotChanged(t *testing.T) {
	buf1 := NewResponseBuffer(nil)
	buf2 := NewResponseBuffer(nil)
	_, req := NewTestRequest("GET", "/")
	DoNothing(buf1, req)
	buf1.WriteTo(buf2)

	if buf1.HasChanged() {
		t.Error("buf1 is changed, but should not be")
	}

	if buf2.HasChanged() {
		t.Error("buf2 is changed, but should not be")
	}

	if !buf1.IsOk() {
		t.Error("buf1 should be ok, but is not")
	}

	if !buf2.IsOk() {
		t.Error("buf2 should be ok, but is not")
	}
}

func TestResponseBufferStatusCreate(t *testing.T) {
	writeCreate := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(201)
	}

	buf := NewResponseBuffer(nil)
	_, req := NewTestRequest("GET", "/")
	writeCreate(buf, req)

	if buf.Code != 201 {
		t.Errorf("Code of buf should be %d but is :%d", 201, buf.Code)
	}

	if !buf.IsOk() {
		t.Error("buf should be ok, but is not")
	}
}
