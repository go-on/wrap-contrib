package wraps

// does something before the request is handled further

import (
	"fmt"
	"net/http"

	"github.com/go-on/wrap"
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
