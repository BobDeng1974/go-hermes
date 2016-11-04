package main

import "time"

// MobileApp represents a mobile app made by a User.
type MobileApp struct {
	ID           int
	Name         string
	CreationDate time.Time // when this mobile app was added to our system
	Metrics      *[]Metric
	User         *User
}
