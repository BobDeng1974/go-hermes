package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

func main() {
	var err error
	db, err = initDb(os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_NAME"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	// close database connection when main() returns
	defer db.Close()

	router := newRouter()

	// Do not run on port 80 as a load balancer will listen on that port.
	log.Fatalln(http.ListenAndServe(":8080", router))
}
