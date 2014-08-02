package wraps

import (
	"fmt"
	"net/http"
)

// String is an utf-8 string that is a http.Handler
// and a wrap.Wrapper
type String string

// ServeHTTP writes the String to the http.ResponseWriter
func (s String) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, s)
}

// Wrap implements the wrap.Wrapper interface
func (s String) Wrap(http.Handler) http.Handler { return s }

func wrtStr(s string, contentType string, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", contentType)
	fmt.Fprint(rw, s)
}

// TextString is an utf-8 string that is a http.Handler
// and a wrap.Wrapper
type TextString string

// ServeHTTP writes the TextString to the http.ResponseWriter and sets
// Content-Type header to text/plain; charset=utf-8
func (t TextString) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrtStr(string(t), "text/plain; charset=utf-8", rw)
}

// Wrap implements the wrap.Wrapper interface
func (t TextString) Wrap(http.Handler) http.Handler { return t }

// JSONString is an utf-8 string that is a http.Handler
// and a wrap.Wrapper
type JSONString string

// ServeHTTP writes the JSONString to the http.ResponseWriter and sets
// Content-Type header to application/json; charset=utf-8
func (t JSONString) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrtStr(string(t), "application/json; charset=utf-8", rw)
}

// Wrap implements the wrap.Wrapper interface
func (t JSONString) Wrap(http.Handler) http.Handler { return t }

// CSSString is an utf-8 string that is a http.Handler
// and a wrap.Wrapper
type CSSString string

// ServeHTTP writes the CSSString to the http.ResponseWriter and sets
// Content-Type header to text/css; charset=utf-8
func (t CSSString) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrtStr(string(t), "text/css; charset=utf-8", rw)
}

// Wrap implements the wrap.Wrapper interface
func (t CSSString) Wrap(http.Handler) http.Handler { return t }

// HTMLString is an utf-8 string that is a http.Handler
// and a wrap.Wrapper
type HTMLString string

// ServeHTTP writes the HTMLString to the http.ResponseWriter and sets
// Content-Type header to text/html; charset=utf-8
func (t HTMLString) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrtStr(string(t), "text/html; charset=utf-8", rw)
}

// Wrap implements the wrap.Wrapper interface
func (t HTMLString) Wrap(http.Handler) http.Handler { return t }

// JavaScriptString is a type alias for string that is a http.Handler
// and a wrap.Wrapper
type JavaScriptString string

// ServeHTTP writes the JavaScriptString to the http.ResponseWriter and sets
// Content-Type header to application/javascript; charset=utf-8
func (t JavaScriptString) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrtStr(string(t), "application/javascript; charset=utf-8", rw)
}

// Wrap implements the wrap.Wrapper interface
func (t JavaScriptString) Wrap(http.Handler) http.Handler { return t }
