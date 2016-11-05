package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"time"

	"golang.org/x/crypto/scrypt"
)

// user create request length. This will limit how many data we read from
// request, to avoid attacks when someone might send large amounts of data.
const ucrLength = 100000

// User type represents a user (customer) in our system.
type User struct {
	ID           int          `json:"id"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	Password     string       `json:"password,omitempty"`
	Salt         string       `json:"-"` // do not show salt in json response at all
	CreationDate time.Time    `json:"creationDate"`
	Servers      *[]Server    `json:"servers,omitempty"`
	MobileApps   *[]MobileApp `json:"mobileApps,omitempty"`
}

// userCreate() reads request, validates email, checks if user exists,
// saves user to db, and returns a JSON response.
func userCreate(w http.ResponseWriter, r *http.Request) {
	var user User
	var err error
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ucrLength))
	if err != nil {
		// could not read stream
		log.Fatalln(err)
	}

	if err = r.Body.Close(); err != nil {
		// could not close body
		log.Fatalln(err)
	}

	// could not create user type from provided json
	if err = json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(422) // unprocessable entity
		APIResponse{Error: "Unprocessable entity"}.response(w)

		log.Fatalln(err)
		return
	}

	// email validation
	if _, err = mail.ParseAddress(user.Email); err != nil {
		APIResponse{Error: "Invalid email address"}.response(w)
		return
	}

	// check if there's a user with that username/email already
	exist, err := user.getUser()
	if err != nil {
		log.Println(err)
		w.Header().Set(`Status`, string(http.StatusInternalServerError))
		APIResponse{Error: "Could not check if user exists"}.response(w)
		return
	}

	if exist {
		APIResponse{Error: "User already exists"}.response(w)
		return
	}

	user.encryptPassword() // encrypt password
	// save user to db
	err = user.insert()
	if err != nil {
		log.Println(err)
		w.Header().Set(`Status`, string(http.StatusInternalServerError))
		APIResponse{Error: "Could not insert user"}.response(w)
		return
	}

	// user created!
	w.WriteHeader(http.StatusCreated)
	user.Password = "" // hide user password from response
	APIResponse{Message: "User created successfully!", Metadata: user}.response(w)
}

// getUser() queries database to find out if user already exists based on username and email.
func (u *User) getUser() (bool, error) {
	var id int

	// Prepare statement for reading data
	err := db.QueryRow("SELECT id FROM user WHERE username = ? OR email = ?", u.Username, u.Email).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		// user doesn't exist
		return false, nil
	}

	// user already exists
	return true, nil
}

// encryptPassword() uses scrypt library to encrypt user's password. Salt is generated from rand.Reader.
func (u *User) encryptPassword() {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}

	dk, err := scrypt.Key([]byte(u.Password), salt, 16384, 8, 1, 32)
	if err != nil {
		log.Fatalln(err)
	}

	u.Password = fmt.Sprintf("%x", dk)
	u.Salt = fmt.Sprintf("%s", salt)
}

// insert() saves newly created user in database
func (u *User) insert() error {
	stmt, err := db.Prepare("INSERT INTO user (username, password, salt, email, creationDate) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(u.Username, u.Password, u.Salt, u.Email, time.Now())
	if err != nil {
		return err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(userID)

	return nil
}
