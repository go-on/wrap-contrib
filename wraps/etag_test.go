package wraps

import (
	"fmt"
	"net/http"

	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib-testing/wrapstesting"
	"github.com/go-on/wrap-contrib/helper"
	// "fmt"
	. "launchpad.net/gocheck"
)

type etagSuite struct{}

var _ = Suite(&etagSuite{})

type ctx2 struct {
	http.ResponseWriter
}

func (c *ctx2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("~" + r.URL.Path + "~"))
}

func (c *ctx2) Put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("did put to " + r.URL.Path))
}

func (s *etagSuite) TestETagIfNoneMatch(c *C) {
	r := wrap.New(
		// LOGGER("If-None-Match"),
		IfNoneMatch,
		// LOGGER("ETag"),
		ETag,
		wrap.Handler(wrapstesting.HandlerMethod((*ctx2).ServeHTTP)),
	)

	rw, req := helper.NewTestRequest("GET", "/path")
	r.ServeHTTP(rw, req)
	_et := rw.Header().Get("ETag")

	c.Assert(_et, Not(Equals), "")
	helper.AssertResponse(rw, "~/path~", 200)

	rw, req = helper.NewTestRequest("GET", "/path")
	req.Header.Set("If-None-Match", fmt.Sprintf("%#v", _et))
	r.ServeHTTP(rw, req)
	c.Assert(rw.Header().Get("ETag"), Equals, _et)
	c.Assert(rw.Code, Equals, 304)

	rw, req = helper.NewTestRequest("GET", "/path")
	req.Header.Set("If-None-Match", `"x"`)
	r.ServeHTTP(rw, req)
	c.Assert(rw.Header().Get("ETag"), Equals, _et)
	err := helper.AssertResponse(rw, "~/path~", 200)

	c.Assert(err, Equals, nil)
}

func (s *etagSuite) TestETagIfMatch(c *C) {
	r0 := wrap.New(
		ETag,
		wrap.Handler(&ctx2{}),
		// LOGGER("ETag"),
	)
	r1 := wrap.New(
		IfMatch(r0),
		wrap.Handler(wrapstesting.HandlerMethod((*ctx2).Put)),
		// LOGGER("IfMatch"),
		// LOGGER("PUT"),
	)

	rw, req := helper.NewTestRequest("HEAD", "/path/")
	r0.ServeHTTP(rw, req)
	_et := rw.Header().Get("ETag")

	c.Assert(_et, Not(Equals), "")

	rw, req = helper.NewTestRequest("PUT", "/path/")
	req.Header.Set("If-Match", fmt.Sprintf("%#v", _et))
	r1.ServeHTTP(rw, req)

	err := helper.AssertResponse(rw, "did put to /path/", 200)

	c.Assert(err, Equals, nil)

	rw, req = helper.NewTestRequest("PUT", "/path/")
	req.Header.Set("If-Match", `"x"`)
	r1.ServeHTTP(rw, req)

	err = helper.AssertResponse(rw, "", 412)
	c.Assert(err, Equals, nil)
}
