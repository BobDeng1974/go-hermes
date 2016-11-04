package main

import "github.com/gorilla/mux"

// NewRouter is application router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)
	}

	return router
}
