package main

import (
	// we need this import to be able to use mysql database driver
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// initDb function is responsible for initialising database connection
// and veryfying connection was successful.
func initDB(username, password, dbName string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, dbName))
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	return conn, err
}
