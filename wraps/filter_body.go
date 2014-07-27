package wraps

import (
	"github.com/go-on/method"
	"github.com/go-on/wrap"
	"github.com/go-on/wrap-contrib/helper"

	"net/http"
)

type filterBody struct {
	method method.Method
}

func (f *filterBody) Wrap(next http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !f.method.Is(r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		checked := helper.NewCheckedResponseWriter(w, func(ck *helper.CheckedResponseWriter) bool {
			ck.WriteHeadersTo(w)
			ck.WriteCodeTo(w)
			return false
		})
		next.ServeHTTP(checked, r)
	})
}

// Filter the body for the given method(s)
// to filter mutiple method, use FilterBody(method.PATCH|method.OPTIONS)
func FilterBody(m method.Method) wrap.Wrapper {
	return &filterBody{m}
}
