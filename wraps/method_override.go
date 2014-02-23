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

func (ø methodOverride) ServeHandle(in http.Handler, w http.ResponseWriter, r *http.Request) {
	override := r.Header.Get(overrideHeader)

	if override != "" {
		expectedMethod, accepted := acceptedOverrides[override]
		if !accepted {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `Unsupported value for %s: %#v.
Supported values are PUT, DELETE, PATCH and OPTIONS`, overrideHeader, override)
			return
		}

		if expectedMethod != r.Method {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, `%s with value %s only allowed for %s requests.`,
				overrideHeader, override, expectedMethod)
			return
		}
		// everything went fine, override the method
		r.Header.Del(overrideHeader)
		r.Method = override
	}

	in.ServeHTTP(w, r)

}

func (ø methodOverride) Wrap(in http.Handler) (out http.Handler) {
	return wrap.ServeHandle(ø, in)
}

func MethodOverride() methodOverride {
	return methodOverride{}
}
