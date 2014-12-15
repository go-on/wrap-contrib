package wraps

/*
import (
	"net/http"

	"gopkg.in/go-on/wrap-contrib.v2/helper"
)

type panicOnWrongOrder struct{}

func (p panicOnWrongOrder) Wrap(next http.Handler) http.Handler {
	// NewPanicOnWrongWriteOrder

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		repl := helper.NewPanicOnWrongWriteOrder(rw)
		next.ServeHTTP(repl, req)
	})
}

var PanicOnWrongOrder = panicOnWrongOrder{}

// NewPanicOnWrongWriteOrder
*/
