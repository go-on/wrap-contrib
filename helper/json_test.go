package helper

import (
	"fmt"
	"net/http"
	"strings"

	"testing"
)

func TestJSONResponse(t *testing.T) {
	rec, _ := NewTestRequest("GET", "/")

	var obj = map[string]string{"key": "value"}
	err := JSONResponse(obj, rec)
	if err != nil {
		t.Error(err)
	}
	err = AssertResponse(rec, `{"key":"value"}`, 200)
	if err != nil {
		t.Error(err)
	}

	err = AssertHeader(rec, "Content-Type", "application/json; charset=utf-8")
	if err != nil {
		t.Error(err)
	}
}

func TestJSONResponseError(t *testing.T) {
	var obj = map[string]func(){"key": func() {}}

	rec, _ := NewTestRequest("GET", "/")
	err := JSONResponse(obj, rec)
	if err == nil {
		t.Error("should not be able to build json for a function")
	}
}

func TestJSONRequest(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"key":"value"}`))

	var obj = map[string]string{}
	err := JSONRequest(&obj, req)

	if err != nil {
		t.Error(err)
	}

	val, has := obj["key"]

	if !has {
		t.Errorf("should have scanned \"key\" but has not")
	}

	if val != "value" {
		t.Errorf("should have value \"value\" but has %#v", val)
	}
}

func TestJSONRequestInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"key":"value"`))

	var obj = map[string]string{}
	err := JSONRequest(&obj, req)

	if err == nil {
		t.Error("should return error")
	}
}

func TestJSONRequestInvalidTarget(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"key":"value}"`))

	var obj = map[string]func(){}
	err := JSONRequest(&obj, req)

	if err == nil {
		t.Error("should return error")
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("can't read")
}

func TestJSONRequestReadError(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", errorReader{})

	var obj = map[string]func(){}
	err := JSONRequest(&obj, req)

	if err == nil {
		t.Error("should return error")
	}

	if err.Error() != "can't read" {
		t.Error("wrong error message")
	}
}
