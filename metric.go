package main

import (
	"database/sql"
	"net/http"
)

type metricHandler struct {
	db *sql.DB
}

// adds a metric to time series database
func (mh *metricHandler) metricCreate(w http.ResponseWriter, r *http.Request) {
	// @todo credentials for time series database
	// @todo connection with time series database
}
