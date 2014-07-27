package helper

import (
	"encoding/json"
	"net/http"
)

// JSONResponse mashals the given obj to the given http.ResponseWriter
func JSONResponse(obj interface{}, w http.ResponseWriter) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(obj)
}

// JSONRequest unmarshals the http.Request body to the object of the pointer ptr
func JSONRequest(ptr interface{}, r *http.Request) (err error) {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(ptr)
}
