package main

import "net/http"

// Route type holds information about HTTP Route
type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

// Routes is a collection of Route type
type Routes []Route

var routes = Routes{
	Route{"UserCreate", "POST", "/user/create", UserCreate},
}
