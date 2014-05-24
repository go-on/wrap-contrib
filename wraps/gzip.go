package wraps

// modelled after code of Andrew Gerrand as proposed here https://groups.google.com/forum/#!searchin/golang-nuts/http%2420gzip/golang-nuts/eVnTcMwNVjM/u0a6TQLagnkJ

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type _gzip struct{}

var GZip = _gzip{}

func (g _gzip) Wrap(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			inner.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		/*
			if err != nil {
				http.Error(w, err.String(), http.StatusInternalServerError)
				return
			}
		*/
		defer gz.Close()
		inner.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}