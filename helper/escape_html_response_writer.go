package helper

import (
	"net/http"
	"unicode/utf8"
)

var (
	//similar to http://golang.org/src/pkg/html/escape.go
	ampOrig = []byte(`&`)[0]
	ampRepl = []byte(`&amp;`)

	sgQuoteOrig = []byte(`'`)[0]
	sgQuoteRepl = []byte(`&#39;`)

	dblQuoteOrig = []byte(`"`)[0]
	dblQuoteRepl = []byte(`&#34;`)

	ltQuoteOrig = []byte(`<`)[0]
	ltQuoteRepl = []byte(`&lt;`)

	gtQuoteOrig = []byte(`>`)[0]
	gtQuoteRepl = []byte(`&gt;`)
)

// EscapeHTMLResponseWriter wraps an http.ResponseWriter in order to override
// its Write method so that it escape html special chars while writing
type EscapeHTMLResponseWriter struct {
	http.ResponseWriter
}

// Write writes to the inner *http.ResponseWriter escaping html special chars on the fly
// Since there is nothing useful to do with the number of bytes written returned from
// the inner responsewriter, the returned int is always 0. Since there is nothing useful to do
// in case of a failed write to the response writer, writing errors are silently dropped.
// the method is modelled after EscapeText from encoding/xml
func (rw *EscapeHTMLResponseWriter) Write(b []byte) (num int, err error) {
	var esc []byte
	n := len(b)
	last := 0

	for i := 0; i < n; {
		r, width := utf8.DecodeRune(b[i:])
		i += width
		switch r {
		case '&':
			esc = ampRepl
		case '\'':
			esc = sgQuoteRepl
		case '"':
			esc = dblQuoteRepl
		case '<':
			esc = ltQuoteRepl
		case '>':
			esc = gtQuoteRepl
		default:
			continue
		}

		rw.ResponseWriter.Write(b[last : i-width])
		rw.ResponseWriter.Write(esc)
		last = i
	}

	rw.ResponseWriter.Write(b[last:])
	return
}
