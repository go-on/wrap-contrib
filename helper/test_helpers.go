package helper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Write is a string that is a simple http.Handler that writes itself to the http.ResponseWriter
type Write string

// ServeHTTP writes the string to the http.ResponseWriter, sets the Content-Length and
// the Content-Type to text/plain
func (ww Write) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(ww)))
	fmt.Fprintf(w, string(ww))
}

// NewTestRequest creates a testing request and a response recorder
func NewTestRequest(method, path string) (recorder *httptest.ResponseRecorder, request *http.Request) {
	request, _ = http.NewRequest(method, path, nil)
	recorder = httptest.NewRecorder()
	return
}

// AssertResponse checks if the given body and code matches the response recorder.
// If one of them do not match, an error is returned
func AssertResponse(rec *httptest.ResponseRecorder, body string, code int) error {
	trimmed := strings.TrimSpace(string(rec.Body.Bytes()))
	if trimmed != body {
		return fmt.Errorf("wrong response body, should be: %#v, but is: %#v", body, trimmed)
	}

	if rec.Code != code {
		return fmt.Errorf("wrong response status code, should be: %d, but is: %d", code, rec.Code)
	}
	return nil
}

// AssertHeader checks, if the ResponseRecorder has a header of key with the value val.
// If it has not, an error is returned.
func AssertHeader(rec *httptest.ResponseRecorder, key, val string) error {
	v := rec.Header().Get(key)
	if v != val {
		return fmt.Errorf("wrong response header, should be: %#v, but is: %#v", val, v)
	}
	return nil
}

// NotFound writes a status code of 404 and the body "not found" to the ResponseWriter
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(`not found`))
}

// DoNothing is a http.Handler that does nothing
func DoNothing(rw http.ResponseWriter, req *http.Request) {}
