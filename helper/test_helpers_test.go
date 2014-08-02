package helper

import "testing"

func TestAssertResponseNoError(t *testing.T) {
	rec, _ := NewTestRequest("GET", "/")
	rec.WriteHeader(201)
	rec.Write([]byte("hi"))

	err := AssertResponse(rec, "hi", 201)
	if err != nil {
		t.Error(err)
	}
}

func TestAssertResponseErrorBody(t *testing.T) {
	rec, _ := NewTestRequest("GET", "/")
	rec.WriteHeader(201)
	rec.Write([]byte("hi"))

	err := AssertResponse(rec, "ho", 201)
	if err == nil {
		t.Error("should return error for wrong body")
	}
}

func TestAssertResponseErrorStatus(t *testing.T) {
	rec, _ := NewTestRequest("GET", "/")
	rec.WriteHeader(201)
	rec.Write([]byte("hi"))

	err := AssertResponse(rec, "hi", 200)
	if err == nil {
		t.Error("should return error for wrong status")
	}
}

func TestAssertHeaderNoError(t *testing.T) {
	rec, req := NewTestRequest("GET", "/")
	write("hi").ServeHTTP(rec, req)

	err := AssertHeader(rec, "Content-Type", "text/plain")

	if err != nil {
		t.Error(err)
	}
}

func TestAssertHeaderError(t *testing.T) {
	rec, req := NewTestRequest("GET", "/")
	write("hi").ServeHTTP(rec, req)

	err := AssertHeader(rec, "Content-Type", "text-plain")

	if err == nil {
		t.Error("should return error for wrong Content-Type")
	}
}

func TestWriteError(t *testing.T) {
	rec, req := NewTestRequest("GET", "/")
	WriteError(rec, req)

	err := AssertResponse(rec, "500 server error", 500)

	if err != nil {
		t.Error("should return server error (500)")
	}
}

/*
func ExamplePrintBody(t *testing.T) {
	rec, req := NewTestRequest("GET", "/")
	write("here is the body").ServeHTTP(rec, req)
	PrintBody(rec)

	// Output
	// here is the body
}
*/
