package wraps

import (
	"github.com/go-on/method"
	"github.com/go-on/wrap"

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

		checked := wrap.NewPeek(w, func(ck *wrap.Peek) bool {
			ck.FlushHeaders()
			ck.FlushCode()
			return false
		})
		next.ServeHTTP(checked, r)

		checked.FlushMissing()
	})
}

// Filter the body for the given method
func FilterBody(m method.Method) wrap.Wrapper {
	return &filterBody{m}
}
