package wraps

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-on/method"
	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"
)

var etagMethods = method.GET | method.HEAD

type etag struct{}

var ETag = etag{}

func (e etag) ServeHandle(in http.Handler, w http.ResponseWriter, r *http.Request) {
	buf := helper.NewResponseBuffer(w)
	// fmt.Printf("in is %T\n", in)
	in.ServeHTTP(buf, r)
	m, _ := method.StringToMethod[r.Method]
	b := buf.Body()

	//	if fake.IsOk() && fake.WHeader != 206 {

	// fmt.Println("Status code", buf.Code)
	// set etag only for status code 200
	if buf.Code == 0 || buf.Code == 200 {
		var etag string
		h := md5.New()
		_, err := io.Copy(h, &buf.Buffer)
		if err == nil {
			// _, err = io.WriteString(h, r.URL.Path)
			if err == nil {
				etag = fmt.Sprintf("%x", h.Sum(nil))
			}
		}
		if etag != "" && etagMethods&m != 0 {
			// fmt.Printf("setting ETag to: %#v for  method %s\n", etag, m.String())
			buf.Header().Set("ETag", etag)
		}
	}

	buf.WriteHeadersTo(w)
	if buf.Code != 0 {
		w.WriteHeader(buf.Code)
	}
	if m != method.HEAD {
		w.Write(b)
	}
}

func (e etag) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(e, inner)
}

type ifNoneMatch struct{}

var IfNoneMatch = ifNoneMatch{}

// see http://www.freesoft.org/CIE/RFC/2068/187.htm
func (i ifNoneMatch) ServeHandle(in http.Handler, w http.ResponseWriter, r *http.Request) {
	ifnone := r.Header.Get("If-None-Match")
	ifnone = strings.Trim(ifnone, `"`)
	// proceed as normal
	if ifnone == "" {
		in.ServeHTTP(w, r)
		return
	}

	ver, _ := method.StringToMethod[r.Method]
	// return 412 for method other than GET and HEAD
	if etagMethods&ver == 0 {
		w.WriteHeader(412) // precondition failed
		return
	}

	buf := helper.NewResponseBuffer(w)
	in.ServeHTTP(buf, r)

	// non 2xx returns should ignire If-None-Match
	if !buf.IsOk() {
		buf.WriteHeadersTo(w)
		if buf.Code != 0 {
			w.WriteHeader(buf.Code)
		}
		buf.WriteTo(w)
		return
	}

	buf.WriteHeadersTo(w)

	// if we have an etag and If-None-Match == * or if the If-None-Match header matches
	// do nothing, but only return the ETag and 304 status
	// fmt.Printf("fake headers: %#v\n", fake.Header())
	etag := buf.Header().Get("ETag")

	// fmt.Println("ifnonematch", ifnone, "etag", etag)
	if (ifnone == "*" && etag != "") || ifnone == etag {
		w.WriteHeader(304)
		return
	}

	// here we return the ressource normally
	if buf.Code != 0 {
		buf.WriteCodeTo(w)
	}
	buf.WriteTo(w)
}

func (i ifNoneMatch) Wrap(inner http.Handler) http.Handler {
	return wrap.ServeHandle(i, inner)
}

type ifMatch struct {
	http.Handler
}

var ifMatchMethods = method.GET | method.PUT | method.DELETE | method.PATCH

// Wo macht If-Match sinn?
// Get, Put, Patch, Delete, alle mit ressource

// see http://stackoverflow.com/questions/2157124/http-if-none-match-vs-if-match
func (i *ifMatch) Wrap(in http.Handler) (out http.Handler) {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ifmatch := r.Header.Get("If-Match")

		// proceed as normal
		ifmatch = strings.Trim(ifmatch, `"`)
		// fmt.Printf("ifmatch: %#v\n", ifmatch)

		if ifmatch == "" || ifmatch == "*" {
			in.ServeHTTP(w, r)
			return
		}

		r.Header.Del("If-Match")

		m, _ := method.StringToMethod[r.Method]
		// return 412 for method other than GET and PUT and DELETE and PATCH
		if ifMatchMethods&m == 0 {
			// fmt.Println("precondition failed, method", m.String())
			w.WriteHeader(412) // precondition failed
			return
		}

		// fmt.Println("make head request")
		buf := helper.NewResponseBuffer(w)
		//	r.Method = "HEAD"
		r2, _ := http.NewRequest("HEAD", r.URL.Path, nil)
		i.ServeHTTP(buf, r2)

		// fmt.Println("code", buf.Code)
		var etag string
		if buf.IsOk() {
			etag = buf.Header().Get("ETag")
			// fmt.Printf("etag %#v ifmatch %#v", etag, ifmatch)
		}

		// fmt.Printf("returned to ifmatch: %#v\n", ifmatch)
		// fmt.Printf("fake headers: %#v\n", fake.Header())
		// fmt.Printf("real headers: %#v\n", w.Header())

		if etag == "" || etag != ifmatch {
			if etag != "" && etagMethods&m != 0 {
				w.Header().Set("ETag", etag)
			}
			// fmt.Println("precondition failed, etag does not match")
			w.WriteHeader(412) // precondition failed
			return
		}

		in.ServeHTTP(w, r)
	})
}

// the given handler will receive a head request for the same path
// and may set an etag in the response
// if it does so, the etag will be compared to the IfMatch header
func IfMatch(h http.Handler) wrap.Wrapper {
	return &ifMatch{h}
}
