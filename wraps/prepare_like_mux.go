package wraps

import (
	"net/http"
	"path"

	"gopkg.in/go-on/wrap.v2"
)

// this is taken from pkg/net

// Return the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

type prepareLikeMux struct{}

// PrepareLikeMux prepares a request the same way that net/http ServeMux does
func PrepareLikeMux() wrap.Wrapper {
	return prepareLikeMux{}
}

func (p prepareLikeMux) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		if rq.RequestURI == "*" {
			if rq.ProtoAtLeast(1, 1) {
				w.Header().Set("Connection", "close")
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// taken from net/http ServeMux and modified
		if rq.Method != "CONNECT" {
			if p := cleanPath(rq.URL.Path); p != rq.URL.Path {
				url := *rq.URL
				url.Path = p
				http.RedirectHandler(url.String(), http.StatusMovedPermanently).ServeHTTP(w, rq)
				return
			}
		}
		next.ServeHTTP(w, rq)
	})
}
