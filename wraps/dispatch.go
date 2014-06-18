package wraps

import (
	"net/http"

	"github.com/go-on/wrap"
)

type Dispatcher interface {
	Dispatch(*http.Request) http.Handler
}

type StructDispatcher interface {
	New(*http.Request) http.HandlerFunc
}

type FuncDispatcher func(*http.Request) http.HandlerFunc

func (fd FuncDispatcher) Dispatch(r *http.Request) http.Handler {
	h := fd(r)
	if h == nil {
		return nil
	}
	return h
}

type DispatchFunc func(*http.Request) http.Handler

func (df DispatchFunc) Dispatch(r *http.Request) http.Handler {
	return df(r)
}

func (df DispatchFunc) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	disp := df(req)
	if disp == nil {
		inner.ServeHTTP(wr, req)
		return
	}
	disp.ServeHTTP(wr, req)
}

func (df DispatchFunc) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	disp := df(req)
	if disp == nil {
		return
	}
	disp.ServeHTTP(wr, req)
}

func (df DispatchFunc) Wrap(inner http.Handler) (out http.Handler) {
	return wrap.ServeHandle(df, inner)
}

func Dispatch(d Dispatcher) DispatchFunc {
	return DispatchFunc(d.Dispatch)
}

func FuncDispatch(fn FuncDispatcher) DispatchFunc {
	return Dispatch(fn)
}

func StructDispatch(stru StructDispatcher) DispatchFunc {
	return FuncDispatch(stru.New)
}

type matchHandler struct {
	Matcher
	http.Handler
}

type MethodHandler struct {
	GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD http.Handler
}

func (mh *MethodHandler) handlerForMethod(method string) http.Handler {
	switch method {
	case "GET":
		return mh.GET
	case "HEAD":
		return mh.HEAD
	case "PUT":
		return mh.PUT
	case "PATCH":
		return mh.PATCH
	case "DELETE":
		return mh.DELETE
	case "OPTIONS":
		return mh.OPTIONS
	case "POST":
		return mh.POST
	default:
		return nil
	}
}

func (mh *MethodHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	h := mh.handlerForMethod(req.Method)
	if h != nil {
		h.ServeHTTP(wr, req)
	}
}

func (mh *MethodHandler) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	h := mh.handlerForMethod(req.Method)
	if h == nil {
		inner.ServeHTTP(wr, req)
		return
	}
	h.ServeHTTP(wr, req)
}

func (mh *MethodHandler) Wrap(inner http.Handler) (out http.Handler) {
	return wrap.ServeHandle(mh, inner)
}

func POSTHandler(path string, handler http.Handler) *matchHandler {
	return &matchHandler{And(MatchPath(path), MatchMethod("POST")), handler}
}

func GETHandler(path string, handler http.Handler) *matchHandler {
	return &matchHandler{And(MatchPath(path), MatchMethod("GET")), handler}
}

func PUTHandler(path string, handler http.Handler) *matchHandler {
	return &matchHandler{And(MatchPath(path), MatchMethod("PUT")), handler}
}

func DELETEHandler(path string, handler http.Handler) *matchHandler {
	return &matchHandler{And(MatchPath(path), MatchMethod("DELETE")), handler}
}

func PATCHHandler(path string, handler http.Handler) *matchHandler {
	return &matchHandler{And(MatchPath(path), MatchMethod("PATCH")), handler}
}

func OPTIONSHandler(path string, handler http.Handler) *matchHandler {
	return &matchHandler{And(MatchPath(path), MatchMethod("OPTIONS")), handler}
}

func HEADHandler(path string, handler http.Handler) *matchHandler {
	return &matchHandler{And(MatchPath(path), MatchMethod("HEAD")), handler}
}

func (mh *matchHandler) ServeHandle(inner http.Handler, wr http.ResponseWriter, req *http.Request) {
	if mh.Match(req) {
		mh.Handler.ServeHTTP(wr, req)
		return
	}
	inner.ServeHTTP(wr, req)
}

func (mh *matchHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if mh.Match(req) {
		mh.Handler.ServeHTTP(wr, req)
	}
}

func (mh *matchHandler) Wrap(inner http.Handler) (out http.Handler) {
	return wrap.ServeHandle(mh, inner)
}

type dispatchMap []matchHandler

// first match counts
func (dm dispatchMap) Dispatch(r *http.Request) http.Handler {
	for _, mh := range dm {
		if mh.Match(r) {
			return mh.Handler
		}
	}
	return nil
}

// data should be pairs of Matcher and http.Handler
func Map(data ...interface{}) DispatchFunc {
	m := dispatchMap{}
	for i := 0; i < len(data); i += 2 {
		m = append(m, matchHandler{data[i].(Matcher), data[i+1].(http.Handler)})
	}
	//return &dispatch{m}
	return Dispatch(m)
}
