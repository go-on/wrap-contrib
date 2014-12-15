package wraps

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"net/http"
	"strings"

	"gopkg.in/go-on/method.v1"
	"gopkg.in/go-on/wrap.v2"
)

type etag struct{}

// ETag buffers the body writes of the next handler and calculates a md5 Hash based on the Content-Type + Body
// combination and sets it as etag in the response header.
// It does so only for GET and HEAD requests. For GET requests the buffered body is flushed to the underlying response writer.
var ETag = etag{}

type etaggedWriter struct {
	*wrap.Peek
	h       hash.Hash
	buf     *bytes.Buffer
	gotData bool
}

func (et *etaggedWriter) Write(b []byte) (num int, err error) {
	if et.buf != nil {
		num, err = et.buf.Write(b)
	}
	if et.IsOk() {
		if !et.gotData {
			ctype := et.Header().Get("Content-Type")
			if ctype != "" {
				et.h.Write([]byte(ctype))
			}
		}
		et.h.Write(b)
	}
	et.gotData = true
	return 0, nil
}

func (e etag) ServeHTTPNext(next http.Handler, w http.ResponseWriter, r *http.Request) {
	m := method.Method(r.Method)

	if !m.MayHaveEtag() {
		next.ServeHTTP(w, r)
		return
	}
	et := &etaggedWriter{h: md5.New(), Peek: wrap.NewPeek(w, nil)}

	// cache for non HEAD methods
	if m != method.HEAD {
		et.buf = bytes.NewBuffer(nil)
	}

	next.ServeHTTP(et, r)
	if et.IsOk() && et.gotData {
		et.Header().Set("ETag", fmt.Sprintf("%x", et.h.Sum(nil)))
	}
	et.FlushMissing()

	if et.buf != nil {
		et.buf.WriteTo(w)
	}
}

// Wrap implements the wrap.Wrapper interface
func (e etag) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(e).Wrap(next)
}

type ifNoneMatch struct{}

// IfNoneMatch only acts upon HEAD and GET requests that have a If-None-Match header set.
// It then runs the next handler and checks if an Etag header is set and if it matches the
// If-None-Match. It it does, the status code 304 (not modified) is sent (without body).
// Otherwise the response is flushed as is.
// The combination of Etag and IfNoneMatch can be used to trigger effective client side caching.
// But IfNoneMatch may also be used with a custom handler that sets the ETag header.
var IfNoneMatch = ifNoneMatch{}

// see http://www.freesoft.org/CIE/RFC/2068/187.htm
func (i ifNoneMatch) ServeHTTPNext(next http.Handler, w http.ResponseWriter, r *http.Request) {
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

	checked := wrap.NewPeek(w, func(ck *wrap.Peek) bool {
		ck.FlushHeaders()
		// non 2xx returns should ignore If-None-Match
		if !ck.IsOk() {
			ck.FlushCode()
			return true
		}

		// if we have an etag and If-None-Match == * or if the If-None-Match header matches
		// do not write the body, but only return the ETag and 304 (not modified) status
		etag := w.Header().Get("ETag")
		if (ifnone == "*" && etag != "") || ifnone == etag {
			w.WriteHeader(http.StatusNotModified)
			w.Write([]byte("\n"))
			return false
		}
		ck.FlushCode()
		return true
	})

	next.ServeHTTP(checked, r)

	checked.FlushMissing()
}

// Wrap implements the wrap.Wrapper interface.
func (i ifNoneMatch) Wrap(next http.Handler) http.Handler {
	return wrap.NextHandler(i).Wrap(next)
}

type ifMatch struct {
	http.Handler
}

// Wrap implements the wrap.Wrapper interface.
//
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

		checkedHead := wrap.NewPeek(w, nil)

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

// IfMatch only acts upon GET, PUT, DELETE or PATCH requests, that have the
// If-Match header set. It then issues a HEAD request to the given handler
// in order to get the etag from the header and compares the etag to the one
// given via the If-Match header. If it matches, the next handler is called,
// otherwise status 412 (precondition failed) is returned
// IfMatch may be used in combination with a middleware stack that uses ETag or
// any custom handler that sets the ETag header.
func IfMatch(h http.Handler) wrap.Wrapper {
	return &ifMatch{h}
}
