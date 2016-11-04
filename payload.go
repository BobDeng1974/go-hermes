package main

// Payload holds the metrics we receive in a request.
type Payload struct {
	User        int    `json:"user_id"`       // payload is specific to a user
	Server      int    `json:"server_id"`     // payload is specific to a server or mobile app
	MobileApp   int    `json:"mobile_app_id"` // payload is specific to a server or mobile app
	MetricID    int    `json:"metric_id"`     // the collected metric from server or mobile app
	MetricValue string `json:"value"`
}
