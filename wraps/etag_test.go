package wraps

import (
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/go-on/wrap.v2"
	"gopkg.in/go-on/wrap-contrib.v2/helper"
	// "fmt"
	//. "launchpad.net/gocheck"
)

/*
type etagSuite struct{}

var _ = Suite(&etagSuite{})
*/
type ctx2 struct {
	http.ResponseWriter
}

func (c *ctx2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	TextContentType.SetContentType(w)
	w.Write([]byte("~" + r.URL.Path + "~"))
}

func (c *ctx2) ServeHTTP2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(302)
	w.Write([]byte("~" + r.URL.Path + "~"))
}

func (c *ctx2) ServeHTTP3(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("~" + r.URL.Path + "~"))
}

func (c *ctx2) Put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("did put to " + r.URL.Path))
}

//func (s *etagSuite) TestETagIfNoneMatch(c *C) {
func TestETagIfNoneMatch(t *testing.T) {
	c := &ctx2{}
	r := wrap.New(
		// LOGGER("If-None-Match"),
		IfNoneMatch,
		// LOGGER("ETag"),
		ETag,
		wrap.HandlerFunc(c.ServeHTTP),
	)

	r1 := wrap.New(
		// LOGGER("If-None-Match"),
		IfNoneMatch,
		// LOGGER("ETag"),
		ETag,
		wrap.HandlerFunc(c.ServeHTTP2),
	)

	r2 := wrap.New(
		// LOGGER("If-None-Match"),
		IfNoneMatch,
		// LOGGER("ETag"),
		ETag,
		wrap.HandlerFunc(c.ServeHTTP3),
	)

	rw, req := helper.NewTestRequest("GET", "/path")
	r.ServeHTTP(rw, req)
	_et := rw.Header().Get("ETag")

	// c.Assert(_et, Not(Equals), "")
	if _et == "" {
		t.Errorf("etag is empty")
	}

	err := helper.AssertResponse(rw, "~/path~", 200)

	if err != nil {
		t.Errorf(err.Error())
	}

	rw, req = helper.NewTestRequest("GET", "/path")
	req.Header.Set("If-None-Match", fmt.Sprintf("%#v", _et))
	r.ServeHTTP(rw, req)

	if rw.Header().Get("ETag") != _et {
		t.Errorf("etag should be %#v but is %#v", _et, rw.Header().Get("ETag"))
	}

	// c.Assert(rw.Header().Get("ETag"), Equals, _et)

	if rw.Code != 304 {
		t.Errorf("code must be 304, but is %v", rw.Code)
	}
	// c.Assert(rw.Code, Equals, 304)

	rw, req = helper.NewTestRequest("PATCH", "/path")
	req.Header.Set("If-None-Match", fmt.Sprintf("%#v", _et))
	r.ServeHTTP(rw, req)

	// c.Assert(rw.Header().Get("ETag"), Equals, _et)

	if rw.Code != 412 {
		t.Errorf("code must be 412, but is %v", rw.Code)
	}

	rw, req = helper.NewTestRequest("GET", "/path")
	req.Header.Set("If-None-Match", fmt.Sprintf("%#v", _et))
	r1.ServeHTTP(rw, req)

	// c.Assert(rw.Header().Get("ETag"), Equals, _et)

	err = helper.AssertResponse(rw, "~/path~", 302)

	if err != nil {
		t.Errorf(err.Error())
	}

	rw, req = helper.NewTestRequest("GET", "/path")
	req.Header.Set("If-None-Match", `"x"`)
	r.ServeHTTP(rw, req)

	if rw.Header().Get("ETag") != _et {
		t.Errorf("etag of If-None-Match should be %#v but is %#v", _et, rw.Header().Get("ETag"))
	}

	// c.Assert(rw.Header().Get("ETag"), Equals, _et)
	err = helper.AssertResponse(rw, "~/path~", 200)

	//c.Assert(err, Equals, nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	rw, req = helper.NewTestRequest("GET", "/path")
	req.Header.Set("If-None-Match", `"x"`)
	r2.ServeHTTP(rw, req)
	err = helper.AssertResponse(rw, "~/path~", http.StatusCreated)

	//c.Assert(err, Equals, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
}

//func (s *etagSuite) TestETagIfMatch(c *C) {
func TestETagIfMatch(t *testing.T) {
	c := &ctx2{}
	r0 := wrap.New(
		ETag,
		wrap.Handler(c),
		// LOGGER("ETag"),
	)
	r1 := wrap.New(
		IfMatch(r0),
		wrap.HandlerFunc(c.Put),
		// LOGGER("IfMatch"),
		// LOGGER("PUT"),
	)

	/*
		r3 := wrap.New(
			IfMatch(r2),
			wrap.Handler(wrapstesting.HandlerMethod((*ctx2).Put)),
			// LOGGER("IfMatch"),
			// LOGGER("PUT"),
		)
	*/
	rw, req := helper.NewTestRequest("HEAD", "/path/")
	r0.ServeHTTP(rw, req)
	_et := rw.Header().Get("ETag")

	if _et == "" {
		t.Errorf("etag is empty")
	}

	// c.Assert(_et, Not(Equals), "")
	rw, req = helper.NewTestRequest("PUT", "/path/")
	r1.ServeHTTP(rw, req)
	err := helper.AssertResponse(rw, "did put to /path/", 200)

	// c.Assert(err, Equals, nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	rw, req = helper.NewTestRequest("PATCH", "/path/")
	req.Header.Set("If-Match", fmt.Sprintf("%#v", _et))
	r1.ServeHTTP(rw, req)

	err = helper.AssertResponse(rw, "did put to /path/", 200)

	// c.Assert(err, Equals, nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	rw, req = helper.NewTestRequest("POST", "/path/")
	req.Header.Set("If-Match", fmt.Sprintf("%#v", _et))
	r1.ServeHTTP(rw, req)

	err = helper.AssertResponse(rw, "", 412)

	// c.Assert(err, Equals, nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	rw, req = helper.NewTestRequest("GET", "/path/")
	req.Header.Set("If-Match", fmt.Sprintf("%#v", _et+"x"))
	r1.ServeHTTP(rw, req)

	err = helper.AssertResponse(rw, "", 412)
	et_ := rw.Header().Get("ETag")

	if et_ != _et {
		t.Errorf("etag should be %#v but is %#v", _et, et_)
	}

	// c.Assert(err, Equals, nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	rw, req = helper.NewTestRequest("PUT", "/path/")
	req.Header.Set("If-Match", `"x"`)
	r1.ServeHTTP(rw, req)

	err = helper.AssertResponse(rw, "", 412)

	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestETag(t *testing.T) {
	c := &ctx2{}
	r2 := wrap.New(
		ETag,
		wrap.HandlerFunc(c.ServeHTTP2),
		// LOGGER("ETag"),
	)
	rw, req := helper.NewTestRequest("GET", "/path/")
	r2.ServeHTTP(rw, req)
	_et2 := rw.Header().Get("ETag")

	err := helper.AssertResponse(rw, "~/path/~", 302)

	if err != nil {
		t.Errorf(err.Error())
	}

	if _et2 != "" {
		t.Errorf("wrong etag: %#v should be %#v", _et2, "")
	}
}

func TestNoETag(t *testing.T) {
	c := &ctx2{}
	r2 := wrap.New(
		ETag,
		wrap.HandlerFunc(c.ServeHTTP2),
		// LOGGER("ETag"),
	)
	rw, req := helper.NewTestRequest("POST", "/path/")
	r2.ServeHTTP(rw, req)
	_et2 := rw.Header().Get("ETag")

	err := helper.AssertResponse(rw, "~/path/~", 302)

	if err != nil {
		t.Errorf(err.Error())
	}

	if _et2 != "" {
		t.Errorf("wrong etag: %#v should be %#v", _et2, "")
	}
}
