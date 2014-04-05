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

func (f *filterBody) Wrap(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ver, _ := method.StringToMethod[r.Method]
		if f.method&ver == 0 {
			in.ServeHTTP(w, r)
			return
		}

		buf := helper.NewResponseBuffer(w)
		in.ServeHTTP(buf, r)
		buf.WriteHeadersTo(w)

		if buf.Code != 0 {
			w.WriteHeader(buf.Code)
		}
	})
}

// Filter the body for the given method(s)
// to filter mutiple method, use FilterBody(method.PATCH|method.OPTIONS)
func FilterBody(m method.Method) wrap.Wrapper {
	return &filterBody{m}
}
