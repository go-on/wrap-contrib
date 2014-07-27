package wraps

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"

	"github.com/go-on/method"
	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
)

// var etagMethods = method.GET | method.HEAD

type etag struct{}

var ETag = etag{}

type etaggedWriter struct {
	*helper.CheckedResponseWriter
	h   hash.Hash
	buf *bytes.Buffer
}

func (et *etaggedWriter) Write(b []byte) (num int, err error) {
	if et.buf != nil {
		num, err = et.buf.Write(b)
	}
	if et.IsOk() {
		et.h.Write(b)
	}
	return 0, io.EOF
}

func (e etag) ServeHandle(next http.Handler, w http.ResponseWriter, r *http.Request) {
	m := method.Method(r.Method)

	if !m.MayHaveEtag() {
		next.ServeHTTP(w, r)
		return
	}
	et := &etaggedWriter{h: md5.New(), CheckedResponseWriter: helper.NewCheckedResponseWriter(w, nil)}

	// cache for non HEAD methods
	if m != method.HEAD {
		et.buf = bytes.NewBuffer(nil)
	}

	next.ServeHTTP(et, r)
	if et.IsOk() {
		et.Header().Set("ETag", fmt.Sprintf("%x", et.h.Sum(nil)))
	}
	et.WriteHeadersTo(w)
	et.WriteCodeTo(w)
	if et.buf != nil {
		et.buf.WriteTo(w)
	}
}

func (e etag) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(e, inner)
}

type ifNoneMatch struct{}

var IfNoneMatch = ifNoneMatch{}

// see http://www.freesoft.org/CIE/RFC/2068/187.htm
func (i ifNoneMatch) ServeHandle(next http.Handler, w http.ResponseWriter, r *http.Request) {
	ifnone := r.Header.Get("If-None-Match")
	ifnone = strings.Trim(ifnone, `"`)
	// proceed as normal
	if ifnone == "" {
		next.ServeHTTP(w, r)
		return
	}

	r.Header.Del("If-None-Match")

	m := method.Method(r.Method)
	// return 412 for method other than GET and HEAD
	if !m.MayHaveEtag() {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	checked := helper.NewCheckedResponseWriter(w, func(ck *helper.CheckedResponseWriter) bool {
		ck.WriteHeadersTo(w)
		// non 2xx returns should ignore If-None-Match
		if !ck.IsOk() {
			ck.WriteCodeTo(w)
			return true
		}

		// if we have an etag and If-None-Match == * or if the If-None-Match header matches
		// do not write the body, but only return the ETag and 304 (not modified) status
		etag := w.Header().Get("ETag")
		if (ifnone == "*" && etag != "") || ifnone == etag {
			w.WriteHeader(http.StatusNotModified)
			return false
		}
		ck.WriteCodeTo(w)
		return true
	})

	next.ServeHTTP(checked, r)
}

func (i ifNoneMatch) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(i, inner)
}

/*
var IfMatchGET = ifMatchGet{}

type ifMatchGet struct{}

// generates the etag and checks it for GET only
func (i *ifMatchGet) Wrap(next http.Handler) (out http.Handler) {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := method.Method(r.Method)

		if m != method.GET {
			next.ServeHTTP(w, r)
		}

		ifmatch := r.Header.Get("If-Match")

		// proceed as normal
		ifmatch = strings.Trim(ifmatch, `"`)

		if ifmatch == "" || ifmatch == "*" {
			next.ServeHTTP(w, r)
			return
		}

		r.Header.Del("If-Match")

		et := &etaggedWriter{h: md5.New(), CheckedResponseWriter: helper.NewCheckedResponseWriter(w, nil)}
		et.buf = bytes.NewBuffer(nil)

		next.ServeHTTP(et, r)
		etag := et.h.Sum(nil)

		if etag != ifmatch {
			w.Header().Set("ETag", etag)
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		et.WriteHeadersTo(w)
		et.WriteCodeTo(w)
		et.buf.WriteTo(w)
	})
}
*/

type ifMatch struct {
	http.Handler
}

// see http://stackoverflow.com/questions/2157124/http-if-none-match-vs-if-match
func (i *ifMatch) Wrap(next http.Handler) (out http.Handler) {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ifmatch := r.Header.Get("If-Match")

		// proceed as normal
		ifmatch = strings.Trim(ifmatch, `"`)

		if ifmatch == "" || ifmatch == "*" {
			next.ServeHTTP(w, r)
			return
		}

		r.Header.Del("If-Match")

		m := method.Method(r.Method)
		// return 412 for method other than GET and PUT and DELETE and PATCH
		if !m.MayHaveIfMatch() {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		checkedHead := helper.NewCheckedResponseWriter(w, nil)
		/*
			checkedHead := helper.NewCheckedResponseWriter(w, func(ck *helper.CheckedResponseWriter) bool {
				return false
			})
		*/

		headReq, _ := http.NewRequest("HEAD", r.URL.Path, nil)
		i.ServeHTTP(checkedHead, headReq)

		var etag string
		if checkedHead.IsOk() {
			etag = checkedHead.Header().Get("ETag")
		}

		if etag == "" || etag != ifmatch {
			if etag != "" && m.MayHaveEtag() {
				w.Header().Set("ETag", etag)
			}
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// the given handler will receive a head request for the same path
// and may set an etag in the response
// if it does so, the etag will be compared to the IfMatch header
func IfMatch(h http.Handler) wrap.Wrapper {
	return &ifMatch{h}
}
