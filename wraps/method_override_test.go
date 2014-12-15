package wraps

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/helper"
)

func methodWrite(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, req.Method)
}

func TestMethodOverride(t *testing.T) {
	r := wrap.New(
		MethodOverride(),
		wrap.HandlerFunc(methodWrite),
	)

	for overrideMethod, requestMethod := range acceptedOverrides {
		rw, req := helper.NewTestRequest(requestMethod, "/")
		req.Header.Set(overrideHeader, overrideMethod)
		r.ServeHTTP(rw, req)
		err := helper.AssertResponse(rw, overrideMethod, 200)
		if err != nil {
			t.Error(err.Error())
		}
		h := req.Header.Get(overrideHeader)
		if h != "" {
			t.Errorf("override header should be removed, but is %#v", h)
		}
	}

	rw, req := helper.NewTestRequest("GET", "/")
	req.Header.Set(overrideHeader, "GET")

	r.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw,
		"Unsupported value for X-HTTP-Method-Override: \"GET\".\nSupported values are PUT, DELETE, PATCH and OPTIONS",
		http.StatusPreconditionFailed)
	if err != nil {
		t.Error(err.Error())
	}

	rw, req = helper.NewTestRequest("GET", "/")
	req.Header.Set(overrideHeader, "PATCH")
	MethodOverride().ServeHTTP(rw, req)

	err = helper.AssertResponse(rw,
		"X-HTTP-Method-Override with value PATCH only allowed for POST requests.",
		http.StatusPreconditionFailed)
	if err != nil {
		t.Error(err.Error())
	}

}

func TestMethodOverrideByFieldDo(t *testing.T) {
	// MethodOverrideByField
	stack := wrap.New(
		MethodOverrideByField("_method"),
		wrap.HandlerFunc(methodWrite),
	)

	vals := url.Values(map[string][]string{})
	vals.Add("_method", "PATCH")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)
	err := helper.AssertResponse(rec, "PATCH", 200)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMethodOverrideByFieldDont1(t *testing.T) {
	// MethodOverrideByField
	stack := wrap.New(
		MethodOverrideByField("_method"),
		wrap.HandlerFunc(methodWrite),
	)

	vals := url.Values(map[string][]string{})
	// vals.Add("_method", "PATCH")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)
	err := helper.AssertResponse(rec, "POST", 200)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMethodOverrideHTTPHandler(t *testing.T) {
	// MethodOverrideByField
	stack := wrap.New(
		Before(MethodOverrideByField("_method")),
		wrap.HandlerFunc(methodWrite),
	)

	vals := url.Values(map[string][]string{})
	// vals.Add("_method", "PATCH")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)
	err := helper.AssertResponse(rec, "POST", 200)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestMethodOverrideByFieldNotAllowed1(t *testing.T) {
	// MethodOverrideByField
	stack := wrap.New(
		MethodOverrideByField("_method"),
		wrap.HandlerFunc(methodWrite),
	)

	vals := url.Values(map[string][]string{})
	vals.Add("_method", "OPTIONS")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)

	if rec.Code != http.StatusPreconditionFailed {
		t.Error("expecting code 412, got: %d", rec.Code)
	}

}

func TestMethodOverrideByFieldNotAllowed2(t *testing.T) {
	// MethodOverrideByField
	stack := wrap.New(
		MethodOverrideByField("_method"),
		wrap.HandlerFunc(methodWrite),
	)

	vals := url.Values(map[string][]string{})
	vals.Add("_method", "MURKS")
	req, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)
	//err := helper.AssertResponse(rec, "_method with value OPTIONS only allowed for GET requests.POST", 412)
	if rec.Code != http.StatusPreconditionFailed {
		t.Error("expecting code 412, got: %d", rec.Code)
	}
}

func TestMethodOverrideByFieldDont2(t *testing.T) {
	// MethodOverrideByField
	stack := wrap.New(
		MethodOverrideByField("_method"),
		wrap.HandlerFunc(methodWrite),
	)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)
	err := helper.AssertResponse(rec, "GET", 200)
	if err != nil {
		t.Error(err.Error())
	}
}
