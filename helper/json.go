package helper

import (
	"encoding/json"
	"net/http"

	"io/ioutil"
)

// JSONResponse mashals the given obj to the given http.ResponseWriter
func JSONResponse(obj interface{}, w http.ResponseWriter) (err error) {
	var b []byte
	b, err = json.Marshal(obj)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
	return
}

// JSONRequest unmarshals the http.Request body to the object of the pointer ptr
func JSONRequest(ptr interface{}, r *http.Request) (err error) {
	var b []byte
	b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	return json.Unmarshal(b, ptr)
}
