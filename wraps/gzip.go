package wraps

// modelled after code of Andrew Gerrand as proposed here https://groups.google.com/forum/#!searchin/golang-nuts/http%2420gzip/golang-nuts/eVnTcMwNVjM/u0a6TQLagnkJ

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/go-on/wrap"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Context(ctxPtr interface{}) {
	w.ResponseWriter.(wrap.Contexter).Context(ctxPtr)
}

func (w *gzipResponseWriter) SetContext(ctxPtr interface{}) {
	w.ResponseWriter.(wrap.Contexter).SetContext(ctxPtr)
}

type _gzip struct{}

var GZip = _gzip{}

func (g _gzip) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		next.ServeHTTP(&gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}
