/*
Package stack provides shortcuts for applications that need just one top level stack
*/
package stack

import (
	"net/http"

	"gopkg.in/go-on/wrap.v2"
)

var context wrap.ContextInjecter

var wrapper []wrap.Wrapper

// New checks inject for validity and sets it as the global top level context
func New(inject wrap.ContextInjecter) {
	if context != nil {
		panic("context already defined")
	}
	if inject == nil {
		context = &basicContext{}
	} else {
		wrap.ValidateContextInjecter(inject)
		context = inject
	}
	wrapper = []wrap.Wrapper{context}
}

// Check may be used in to compose Wrapper with inner middleware stacks like in
// It checks if the wrappers context requirements are fulfilled by the top level context
// stack.Use( wrap.New(stack.Check(wrapper1), stack.Check(wrapper2)) )
func Check(w wrap.Wrapper) wrap.Wrapper {
	if context == nil {
		panic("no context defined, use New")
	}
	wrap.ValidateWrapperContexts(context, w)
	return w
}

// Use adds the given wrappers to the middleware stack after checking if their context
// requirements are fulfilled by the global toplevel context
func Use(w ...wrap.Wrapper) {
	if context == nil {
		panic("no context defined, use New")
	}
	wrap.ValidateWrapperContexts(context, w...)
	wrapper = append(wrapper, w...)
}

// Handler composes the stack and returns a handler for it
func Handler() http.Handler {
	return wrap.New(wrapper...)
}
