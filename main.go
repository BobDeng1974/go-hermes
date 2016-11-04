package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

	log.Fatalln(http.ListenAndServe(":80", router))
}
