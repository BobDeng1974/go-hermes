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
	influxDB "github.com/influxdata/influxdb/client/v2"
	"gopkg.in/mgo.v2"
)

func initMongo(host, username, password string) (*mgo.Session, error) {
	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Username: username,
		Password: password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// initMySQL function is responsible for initialising database connection
// and verifying connection was successful.
func initMySQL(username, password, dbName string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, dbName))
	if err != nil {
		return nil, err
	}

	return conn, conn.Ping()
}

func initInfluxDB(host, username, password string) (influxDB.Client, error) {
	db, err := influxDB.NewHTTPClient(influxDB.HTTPConfig{
		Addr:     host,
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	_, _, err = db.Ping(0)
	if err != nil {
		return nil, err
	}

	return db, nil
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
func loadConfig() (Configuration, error) {
	var c Configuration
	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(f, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func main() {
	c, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// environment variables needed for database connections
	const (
		mysqlUsername = "MYSQL_USERNAME"
		mysqlPassword = "MYSQL_PASSWORD"
		mysqlDBName   = "MYSQL_NAME"
	)

	mysqlDB, err := initMySQL(os.Getenv(mysqlUsername), os.Getenv(mysqlPassword), os.Getenv(mysqlDBName))
	if err != nil {
		log.Fatalln(err)
	}

	influxDBClient, err := initInfluxDB(c.InfluxDBHost, os.Getenv("INFLUX_USER"), os.Getenv("INFLUX_PWD"))
	if err != nil {
		log.Fatalln(err)
	}

	mongoDB, err := initMongo(c.MongoHost, os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASS"))
	if err != nil {
		log.Fatalln(err)
	}

	// close influxDB & MySQL connection when main() returns
	//defer influxDB.Close()
	defer influxDBClient.Close()
	defer mysqlDB.Close()
	defer mongoDB.Close()

	router := newRouter(mysqlDB)

	s := http.Server{
		Addr:         c.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Do not run on port 80 as a load balancer will listen on that port.
	log.Fatalln(s.ListenAndServeTLS(c.Cert, c.CertKey))
}
