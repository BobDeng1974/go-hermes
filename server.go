package main

import "time"

// OS represents an Operating System
type OS struct {
	Name string
}

// Server represents a remote host
type Server struct {
	ID             int
	HostName       string
	User           *User
	LastMetricDate time.Time
	OS             *OS
	Metrics        *[]Metric
	CreationDate   time.Time
}
