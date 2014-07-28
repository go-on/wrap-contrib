package helper

/*
import (
	"bytes"
	"net/http"
)

type panicOnWrongWriteOrder struct {
	http.ResponseWriter
	codeWritten bool
	headerSet   bool
	bodyWritten bool
}
*/

/*
func NewPanicOnWrongWriteOrder(rw http.ResponseWriter) http.ResponseWriter {
	return &panicOnWrongWriteOrder{
		ResponseWriter: rw,
	}
}

// Header returns the http.Header
func (f *panicOnWrongWriteOrder) Header() http.Header {
	if f.codeWritten {
		panic("header set after code written")
	}
	if f.bodyWritten {
		panic("header written after body")
	}
	fmt.Println("writing header")
	f.headerSet = true
	return f.ResponseWriter.Header()
}

func (f *panicOnWrongWriteOrder) RealResonseWriter() http.ResponseWriter {
	start := f.ResponseWriter

	for {
		ck, isChecked := start.(*CheckedResponseWriter)
		if !isChecked {
			return start
		}
		start = ck.ResponseWriter
	}

	return start
}

// WriteHeader writes the status code
func (f *panicOnWrongWriteOrder) WriteHeader(i int) {
	if f.codeWritten {
		panic("code written a second time")
	}
	if f.bodyWritten {
		panic("code written after body")
	}
	fmt.Printf("writing status code: %d, header is: %v (%T)\n", i, f.ResponseWriter.Header(), f.RealResonseWriter())
	f.ResponseWriter.WriteHeader(i)
	f.codeWritten = true
}

func (f *panicOnWrongWriteOrder) Write(b []byte) (int, error) {
	f.bodyWritten = true
	fmt.Printf("writing body (%T): %#v\n", f.RealResonseWriter(), string(b))
	return f.ResponseWriter.Write(b)
}
*/
