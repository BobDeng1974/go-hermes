package main

import "net/http"

// route type holds information about HTTP Route
type route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

// Routes is a collection of Route type
type routes []route

var appRoutes = routes{
	route{"UserCreate", "POST", "/user/create", userCreate},
}
