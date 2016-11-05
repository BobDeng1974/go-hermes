package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIResponse is the HTTP API response our app uses
type APIResponse struct {
	Message  string `json:"message"`
	Metadata interface{}
	Error    string `json:"error,omitempty"`
}

// response() builds and writes the API response.
func (ar APIResponse) response(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// json encode APIResponse in response writer
	if err := json.NewEncoder(w).Encode(ar); err != nil {
		log.Fatalln(err)
	}
}
