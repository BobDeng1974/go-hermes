package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type serverHandler struct {
	session *mgo.Session
	th      tokenHandler
}

const scrLength = 100000 // server create request max length

// serverCreate() decodes request, checks if server already exists. If not creates the server in database.
func (sh *serverHandler) serverCreate(w http.ResponseWriter, r *http.Request) {
	user, err := sh.th.getUserByToken(r.URL.Query().Get("token"))
	if err != nil {
		w.WriteHeader(400) // bad request
		APIResponse{Message: "Invalid token"}.response(w)
		return
	}
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
	server.UserID = user.ID
	// FIXME fails when request has userId value in quotes with error "cannot unmarshal string into Go value of type int".
	if err = json.Unmarshal(body, server); err != nil {
		w.WriteHeader(422) // unprocessable entity
		APIResponse{Message: "Unprocessable entity"}.response(w)

		log.Println(err)
		return
	}

	// check if server already exists based on user ID and hostname
	exist, err := sh.findServer(server.UserID, server.HostName)
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
func (sh *serverHandler) findServer(userID bson.ObjectId, hostname string) (bool, error) {
	var server Server
	c := sh.session.DB("app").C("server")
	err := c.Find(bson.M{"userId": userID, "hostName": hostname}).One(&server)
	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Inserts server in database
func (sh *serverHandler) insert(s *Server) error {
	c := sh.session.DB("app").C("server")
	s.ID = bson.NewObjectId()
	err := c.Insert(bson.M{
		"_id":          s.ID,
		"userId":       s.UserID,
		"hostName":     s.HostName,
		"osName":       s.OS.Name,
		"creationDate": time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}
