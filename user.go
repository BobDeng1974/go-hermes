package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// User type represents a user (customer) in our system.
type User struct {
	ID           int          `json:"id"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	Password     string       `json:"password,omitempty"`
	CreationDate time.Time    `json:"creationDate"`
	Servers      *[]Server    `json:"servers,omitempty"`
	MobileApps   *[]MobileApp `json:"mobileApps,omitempty"`
}

// UserCreate creates a user
func UserCreate(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 100000))
	if err != nil {
		// could not read stream
		log.Fatalln(err)
	}

	if err := r.Body.Close(); err != nil {
		// could not close body
		log.Fatalln(err)
	}

	// could not create user type from provided json
	if err := json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(422) // unprocessable entity
		APIResponse{Error: "Unprocessable entity"}.Response(w)

		log.Fatalln(err)

		return
	}

	user.CreationDate = time.Now()

	// @todo validation
	// @todo encrypt pass, save user to db

	// user created!
	w.WriteHeader(http.StatusCreated)
	user.Password = "" // hide user password from response
	APIResponse{Message: "User created successfully!", Metadata: user}.Response(w)
}
