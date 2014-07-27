package helper

import (
	"io"
	"net/http"

	"github.com/go-on/wrap"
)

type CheckedResponseWriter struct {
	http.ResponseWriter // necessary to allow "Unwrap from wrapstesting"
	Code                int
	changed             bool
	header              http.Header
	writeForbidden      bool
	isChecked           bool
	codeWritten         bool
	headersWritten      bool
	// writtenHeadersAndCode bool
	// Check should return true if the data should be written to the inner ResponseWriter
	// otherwise false
	// Check may check the Code and headers that have been set and to the CheckedResponseWriter
	// and may also decide to transfer them to the inner ResponseWriter or set them directly on
	// the ResponseWriter. Check can be sure to be invoked before the first write to http.ResponseWriter
	check func(*CheckedResponseWriter) bool
}

func NewCheckedResponseWriter(rw http.ResponseWriter, check func(*CheckedResponseWriter) bool) *CheckedResponseWriter {
	return &CheckedResponseWriter{ResponseWriter: rw, check: check, header: make(http.Header)}
}

func (f *CheckedResponseWriter) Context(ctxPtr interface{}) {
	f.ResponseWriter.(wrap.ResponseWriterWithContext).Context(ctxPtr)
}

func (f *CheckedResponseWriter) SetContext(ctxPtr interface{}) {
	f.ResponseWriter.(wrap.ResponseWriterWithContext).SetContext(ctxPtr)
}

// Header returns the http.Header
func (f *CheckedResponseWriter) Header() http.Header {
	f.changed = true
	return f.header
}

// WriteHeader writes the status code
func (f *CheckedResponseWriter) WriteHeader(i int) {
	f.changed = true
	f.Code = i
}

// IsOk returns true if the returned status code is
// not set or in the 2xx range
func (f *CheckedResponseWriter) IsOk() bool {
	if f.Code == 0 {
		return true
	}
	if f.Code >= 200 && f.Code < 300 {
		return true
	}
	return false
}

// Write only writes if writing is allowed
func (f *CheckedResponseWriter) Write(b []byte) (int, error) {
	if !f.isChecked {
		if f.check != nil {
			f.writeForbidden = !f.check(f)
		}
		f.isChecked = true
	}
	if f.writeForbidden {
		return 0, io.EOF
	}
	f.changed = true
	return f.ResponseWriter.Write(b)
}

// Reset set the CheckedResponseWriter to the defaults
func (f *CheckedResponseWriter) Reset() {
	f.Code = 0
	f.changed = false
	f.writeForbidden = false
	f.isChecked = false
	f.codeWritten = false
	f.headersWritten = false
	f.header = make(http.Header)
}

func (f *CheckedResponseWriter) HasChanged() bool {
	return f.changed
}

// WriteCodeTo writes the status code to the responsewriter if it was set
func (f *CheckedResponseWriter) WriteCodeTo(w http.ResponseWriter) {
	if f.codeWritten {
		return
	}
	if f.Code != 0 {
		w.WriteHeader(f.Code)
	}
	f.codeWritten = true
}

// WriteHeadersTo adds the headers to the given ResponseWriter
func (f *CheckedResponseWriter) WriteHeadersTo(w http.ResponseWriter) {
	if f.codeWritten {
		panic("code written before headers")
	}
	if f.headersWritten {
		return
	}
	header := w.Header()
	for k, v := range f.header {
		header.Del(k)
		for _, val := range v {
			header.Add(k, val)
		}
	}
	f.headersWritten = true
}
