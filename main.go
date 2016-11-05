package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// environment variables needed for mysql connection
const (
	mysqlUsername = "MYSQL_USERNAME"
	mysqlPassword = "MYSQL_PASSWORD"
	mysqlDBName   = "MYSQL_NAME"
)

// initDb function is responsible for initialising database connection
// and verifying connection was successful.
func initDB(username, password, dbName string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, dbName))
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	return conn, err
}

// newRouter is application router
func newRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	uh := userHandler{
		db: db,
	}

	router.Methods("POST").
		Path("/user/create").
		Name("UserCreate").
		HandlerFunc(uh.userCreate)

	return router
}

func main() {
	var err error

	// Show a clear message if an environment variable is not set.
	// DB connection will fail without this check, but this check will speed up
	// the debugging process knowing if an environment variable is not set or
	// if the credentials are just wrong.
	if os.Getenv(mysqlUsername) == "" || os.Getenv(mysqlPassword) == "" || os.Getenv(mysqlDBName) == "" {
		log.Fatalln("MySQL database environment variables need to be set")
	}

	db, err := initDB(os.Getenv(mysqlUsername), os.Getenv(mysqlPassword), os.Getenv(mysqlDBName))
	if err != nil {
		log.Fatalln(err.Error())
	}

	// close database connection when main() returns
	defer db.Close()

	router := newRouter(db)

	// Do not run on port 80 as a load balancer will listen on that port.
	log.Fatalln(http.ListenAndServe(":8080", router))
}
