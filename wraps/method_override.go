package wraps

// does something before the request is handled further

import (
	"fmt"
	"net/http"

	"gopkg.in/go-on/method.v1"

	"gopkg.in/go-on/wrap.v2"
)

type methodOverride struct{}

var acceptedOverrides = map[string]string{
	"PATCH":   "POST",
	"OPTIONS": "GET",
	"DELETE":  "POST",
	"PUT":     "POST",
}

var overrideHeader = "X-HTTP-Method-Override"

// returns true if it interrupts (error)
func (ø methodOverride) serveHTTP(w http.ResponseWriter, r *http.Request) (isError bool) {
	override := r.Header.Get(overrideHeader)
	// fmt.Printf("method override called: %v\n", override)

	if override != "" {
		// fmt.Println("override", override)
		expectedMethod, accepted := acceptedOverrides[override]
		if !accepted {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `Unsupported value for %s: %#v.
Supported values are PUT, DELETE, PATCH and OPTIONS`, overrideHeader, override)
			return true
		}

		if expectedMethod != r.Method {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `%s with value %s only allowed for %s requests.`,
				overrideHeader, override, expectedMethod)
			return true
		}
		// everything went fine, override the method
		r.Header.Del(overrideHeader)
		r.Method = override
	}
	return false
}

func (ø methodOverride) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ø.serveHTTP(w, r)
}

func (ø methodOverride) ServeHTTPNext(next http.Handler, w http.ResponseWriter, r *http.Request) {
	if ø.serveHTTP(w, r) {
		return
	}
	next.ServeHTTP(w, r)
}

func (ø methodOverride) Wrap(next http.Handler) (out http.Handler) {
	return wrap.NextHandler(ø).Wrap(next)
}

func MethodOverride() methodOverride {
	return methodOverride{}
}

// MethodOverrideByField overrides the request method by looking for a field that
// contains the target method. It only acts on POST requests and on post bodies.
type MethodOverrideByField string

func (m MethodOverrideByField) serveHTTP(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method.POST.String() {
		return false
	}
	override := r.PostFormValue(string(m))

	// fmt.Printf("override: %#v\n", r.PostForm)

	if override != "" {
		expectedMethod, accepted := acceptedOverrides[override]
		if !accepted {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `Unsupported value for %s: %#v.
Supported values are PUT, DELETE, PATCH and OPTIONS`, m, override)
			return true
		}

		if expectedMethod != r.Method {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `%s with value %s only allowed for %s requests.`,
				m, override, expectedMethod)
			return true
		}
		// everything went fine, override the method
		// r.Header.Del(m)
		r.Form.Del(string(m))
		r.Method = override
	}
	return false
}

func (m MethodOverrideByField) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.serveHTTP(w, r)
}

func (m MethodOverrideByField) Wrap(next http.Handler) (out http.Handler) {
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, req *http.Request) {
		if m.serveHTTP(rw, req) {
			return
		}
		next.ServeHTTP(rw, req)
	}
	return f
}
