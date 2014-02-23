package helper

import (
	"bytes"
	"net/http"
)

// ResponseBuffer is a ResponseWriter that may be used to spy on
// http.Handlers and keep what they have written.
// It may then be written to another (the real) ResponseWriter
type ResponseBuffer struct {
	Buffer  bytes.Buffer
	Code    int
	changed bool
	header  http.Header
}

// Header returns the http.Header
func (f *ResponseBuffer) Header() http.Header {
	f.changed = true
	return f.header
}

// WriteHeader writes the status code
func (f *ResponseBuffer) WriteHeader(i int) { f.changed = true; f.Code = i }

// Write writes to the underlying buffer
func (f *ResponseBuffer) Write(b []byte) (int, error) {
	f.changed = true
	return f.Buffer.Write(b)
}

// Reset set the ResponseBuffer to the defaults
func (f *ResponseBuffer) Reset() {
	f.Buffer.Reset()
	f.Code = 0
	f.changed = false
	f.header = make(http.Header)
}

// WriteTo writes header, body and status code to the given ResponseWriter, if something changed
func (f *ResponseBuffer) WriteTo(wr http.ResponseWriter) {
	if f.HasChanged() {
		f.WriteHeadersTo(wr)
		f.WriteCodeTo(wr)
		wr.Write(f.Buffer.Bytes())
	}
}

// Body returns the body as slice of bytes
func (f *ResponseBuffer) Body() []byte {
	return f.Buffer.Bytes()
}

// BodyString returns the body as string
func (f *ResponseBuffer) BodyString() string {
	return f.Buffer.String()
}

// HasChanged returns true if something has been written to the ResponseBuffer
func (f *ResponseBuffer) HasChanged() bool { return f.changed }

// IsOk returns true if the returned status code is
// not set or in the 2xx range
func (f *ResponseBuffer) IsOk() bool {
	if f.Code == 0 {
		return true
	}
	if f.Code >= 200 && f.Code < 300 {
		return true
	}
	return false
}

// WriteCodeTo writes the status code to the responsewriter if it was set
func (f *ResponseBuffer) WriteCodeTo(w http.ResponseWriter) {
	if f.Code != 0 {
		w.WriteHeader(f.Code)
	}
}

// WriteHeadersTo adds the headers to the given ResponseWriter
func (f *ResponseBuffer) WriteHeadersTo(w http.ResponseWriter) {
	header := w.Header()
	for k, v := range f.header {
		header.Del(k)
		for _, val := range v {
			header.Add(k, val)
		}
	}
}

// NewResponseBuffer creates a new ResponseBuffer
func NewResponseBuffer() (f *ResponseBuffer) {
	f = &ResponseBuffer{}
	f.header = make(http.Header)
	return
}
