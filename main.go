package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

// environment variables needed for mysql connection
const (
	mysqlUsername = "MYSQL_USERNAME"
	mysqlPassword = "MYSQL_PASSWORD"
	mysqlDBName   = "MYSQL_NAME"
)

func main() {
	var err error

	// Show a clear message if an environment variable is not set.
	// DB connection will fail without this check, but this check will speed up
	// the debugging process knowing if an environment variable is not set or
	// if the credentials are just wrong.
	if os.Getenv(mysqlUsername) == "" || os.Getenv(mysqlPassword) == "" || os.Getenv(mysqlDBName) == "" {
		log.Fatalln("MySQL database environment variables need to be set")
	}

	db, err = initDB(os.Getenv(mysqlUsername), os.Getenv(mysqlPassword), os.Getenv(mysqlDBName))
	if err != nil {
		log.Fatalln(err.Error())
	}

	// close database connection when main() returns
	defer db.Close()

	router := newRouter()

	// Do not run on port 80 as a load balancer will listen on that port.
	log.Fatalln(http.ListenAndServe(":8080", router))
}
