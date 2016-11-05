package main

import (
	// we need this import to be able to use mysql database driver
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// initDb function is responsible for initialising database connection
// and veryfying connection was successful.
func initDb(username, password, name string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", username+":"+password+"@/"+name)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	return conn, err
}
