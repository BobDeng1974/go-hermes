package main

import "time"

// OS represents an Operating System
type OS struct {
	Name string
}

// Server represents a remote host
type Server struct {
	ID             int       `json:"id"`
	HostName       string    `json:"hostname"`
	User           *User     `json:"user"`
	LastMetricDate time.Time `json:"lastMetricDate,omitempty"`
	OS             *OS       `json:"os"`
	Metrics        *[]Metric `json:"-"`
	CreationDate   time.Time `json:"creationDate,omitempty"`
}

// A Metric is a measurement that makes sense to User when viewed in dashboard.
type Metric struct {
	ID           int
	Name         string
	Value        string    // metric value
	CreationDate time.Time // when this metric was created on 3rd party host
	ReceivedDate time.Time // when we were notified about this metric's value
	Server       *Server
}

// MobileApp represents a mobile app made by a User.
type MobileApp struct {
	ID           int
	Name         string
	CreationDate time.Time // when this mobile app was added to our system
	Metrics      *[]Metric
	User         *User
}

// Payload holds the metrics we receive in a request.
type Payload struct {
	User        int    `json:"user_id"`       // payload is specific to a user
	Server      int    `json:"server_id"`     // payload is specific to a server or mobile app
	MobileApp   int    `json:"mobile_app_id"` // payload is specific to a server or mobile app
	MetricID    int    `json:"metric_id"`     // the collected metric from server or mobile app
	MetricValue string `json:"value"`
}

// User type represents a user (customer) in our system.
type User struct {
	ID              int          `json:"id"`
	Username        string       `json:"username"`
	Email           string       `json:"email"`
	Password        []byte       `json:"password,omitempty"`
	Salt            []byte       `json:"-"` // do not show salt in json response at all
	CreationDate    time.Time    `json:"creationDate"`
	Servers         *[]Server    `json:"servers,omitempty"`
	MobileApps      *[]MobileApp `json:"mobileApps,omitempty"`
	passwordEncoded bool
}

// Configuration type holds app configuration
type Configuration struct {
	Cert    string `json:"cert"`     // HTTPS certificate
	CertKey string `json:"cert_key"` // HTTPS certificate key
	Port    string `json:"port"`     // port that application runs on

	// influxDB Settings
	InfluxDBHost string `json:"influxDBHost"`

	// mongoDB settings
	MongoHost string `json:"mongoHost"` // mongoDB host
}

// APIToken represents an API Token
type APIToken struct {
	Token string
	User  *User
}
