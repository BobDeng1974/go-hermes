package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type serverHandler struct {
	db *sql.DB
	uh *userHandler
}

const (
	scrLength       = 100000 // server create request max length
	findServerSQL   = "SELECT userID FROM server WHERE userID = ? AND hostname = ? LIMIT 1"
	insertServerSQL = "INSERT INTO server (userId, hostname, os, creationDate) VALUES (?, ?, ?, ?)"
)

// serverCreate() decodes request, checks if server already exists. If not creates the server in database.
func (sh *serverHandler) serverCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, scrLength))
	if err != nil {
		// could not read stream
		log.Fatalln(err)
	}

	err = r.Body.Close()
	if err != nil {
		log.Println(err)
	}

	server := &Server{}
	// FIXME fails when request has userId value in quotes with error "cannot unmarshal string into Go value of type int".
	if err = json.Unmarshal(body, server); err != nil {
		w.WriteHeader(422) // unprocessable entity
		APIResponse{Message: "Unprocessable entity"}.response(w)

		log.Fatalln(err)
		return
	}

	// We want to make sure this request for new server, belongs to an existing user.
	// Therefore, we need to check if user exists.
	exist, err := sh.uh.findByID(server.User.ID)
	if err != nil {
		w.Header().Set(`Status`, string(http.StatusInternalServerError))
		APIResponse{Message: "Could not check if user exists"}.response(w)
		log.Println(err)
		return
	}

	if !exist {
		// user does not exist
		w.Header().Set(`Status`, string(http.StatusNotFound))
		APIResponse{Message: "User not found"}.response(w)
		return
	}

	// check if server already exists based on user ID and hostname
	exist, err = sh.findServer(server.User.ID, server.HostName)
	if err != nil {
		w.Header().Set(`Status`, string(http.StatusInternalServerError))
		APIResponse{Message: "Could not check if server exists"}.response(w)
		log.Println(err)
		return
	}

	if exist {
		// server already exists
		APIResponse{Error: "Server already exists"}.response(w)
		return
	}

	// insert server
	err = sh.insert(server)
	if err != nil {
		w.Header().Set(`Status`, string(http.StatusInternalServerError))
		APIResponse{Message: "Could not insert server"}.response(w)
		log.Println(err)
		return
	}

	// TODO figure out how to exclude User and Server.LastMetric from this API response only
	APIResponse{Message: "Server created successfully!", Metadata: server}.response(w)
}

// Checks if server exists in database based on userID and hostname.
// Returns true if server exists, or false if not.
// Also, returns an error or nil if no error occurred.
func (sh *serverHandler) findServer(userID int, hostname string) (bool, error) {
	err := sh.db.QueryRow(findServerSQL, userID, hostname).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// server does not exist
			return false, nil
		}
		// something went wrong
		return false, err
	}

	// server already exists
	return true, nil
}

// Inserts server in database
func (sh *serverHandler) insert(s *Server) error {
	stmt, err := sh.db.Prepare(insertServerSQL)
	if err != nil {
		return err
	}
	s.CreationDate = time.Now()
	res, err := stmt.Exec(s.User.ID, s.HostName, s.OS.Name, s.CreationDate)
	if err != nil {
		return err
	}

	serverID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(serverID)
	return nil
}
