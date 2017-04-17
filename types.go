package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// OS represents an Operating System
type OS struct {
	Name string
}

// Server represents a remote host
type Server struct {
	ID             bson.ObjectId `json:"id"`
	HostName       string        `json:"hostname"`
	UserID         bson.ObjectId `json:"userId"`
	LastMetricDate time.Time     `json:"lastMetricDate,omitempty"`
	OS             *OS           `json:"os"`
	Metrics        *[]Metric     `json:"-"`
	CreationDate   time.Time     `json:"creationDate,omitempty"`
}

// A Metric is a measurement that makes sense to User when viewed in dashboard.
type Metric struct {
	ID           bson.ObjectId
	Name         string
	Value        string    // metric value
	CreationDate time.Time // when this metric was created on 3rd party host
	ReceivedDate time.Time // when we were notified about this metric's value
	Server       *Server
}

// App represents an (IoT/mobile/web) application made by a User.
type App struct {
	ID           bson.ObjectId
	Name         string
	CreationDate time.Time // when this app was added to our system
	Metrics      *[]Metric
	User         *User
}

// Payload holds the metrics we receive in a request.
type Payload struct {
	User        bson.ObjectId `json:"user_id"`
	Server      bson.ObjectId `json:"server_id"`
	App         bson.ObjectId `json:"app_id"`
	MetricID    bson.ObjectId `json:"metric_id"`
	MetricValue string        `json:"value"`
}

// ACLRole refers to a user role.
type ACLRole string

const (
	// RoleUser is for standard level users.
	RoleUser ACLRole = "ROLE_USER"
	// RoleSuperAdmin is for administrators.
	RoleSuperAdmin ACLRole = "ROLE_SUPER_ADMIN"
)

// User type represents a user (customer) in our system.
type User struct {
	ID              bson.ObjectId `json:"id"`
	Username        string        `json:"username"`
	Email           string        `json:"email"`
	Password        []byte        `json:"password,omitempty"`
	Salt            []byte        `json:"-"` // do not show salt in json response at all
	CreationDate    time.Time     `json:"creationDate"`
	Servers         *[]Server     `json:"servers,omitempty"`
	Apps            *[]App        `json:"apps,omitempty"`
	Roles           []ACLRole     `json:"-"`
	Token           string        `json:"apiToken,omitempty"`
	passwordEncoded bool
	tokenGenerated  bool
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
