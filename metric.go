package main

import "time"

// A Metric is a measurement that makes sense to User when viewed in dashboard.
type Metric struct {
	ID           int
	Name         string
	Value        string    // metric value
	CreationDate time.Time // when this metric was created on 3rd party host
	ReceivedDate time.Time // when we were notified about this metric's value
	Server       *Server
}
