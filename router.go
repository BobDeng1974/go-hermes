package main

import "github.com/gorilla/mux"

// newRouter is application router
func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range appRoutes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)
	}

	return router
}
