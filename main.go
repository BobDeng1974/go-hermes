package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var config Configuration

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

	uh := userHandler{db: db}
	sh := serverHandler{db: db, uh: &uh}

	router.Methods("POST").
		Path("/user/create").
		Name("UserCreate").
		HandlerFunc(uh.userCreate)

	router.Methods("POST").
		Path("/server/create").
		Name("ServerCreate").
		HandlerFunc(sh.serverCreate)

	return router
}

// Reads config.json file, and unmarshals to Configuration struct.
func loadConfig() error {
	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, &config)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Show a clear message if an environment variable is not set.
	// DB connection will fail without this check, but this check will speed up
	// the debugging process knowing if an environment variable is not set or
	// if the credentials are just wrong.
	if os.Getenv(mysqlUsername) == "" || os.Getenv(mysqlPassword) == "" || os.Getenv(mysqlDBName) == "" {
		log.Fatalln("MySQL database environment variables need to be set")
	}

	db, err := initDB(os.Getenv(mysqlUsername), os.Getenv(mysqlPassword), os.Getenv(mysqlDBName))
	if err != nil {
		log.Fatalln(err)
	}

	// close database connection when main() returns
	defer db.Close()

	router := newRouter(db)

	s := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Do not run on port 80 as a load balancer will listen on that port.
	log.Fatalln(s.ListenAndServeTLS(config.Cert, config.CertKey))
}
