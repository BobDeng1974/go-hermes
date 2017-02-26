package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"time"

	"golang.org/x/crypto/scrypt"
)

const (
	insertUserQuery       = "INSERT INTO user (username, password, salt, email, creationDate) VALUES(?, ?, ?, ?, ?)"
	findUserByNameOrEmail = "SELECT id FROM user WHERE username = ? OR email = ?"
	findUserByID          = "SELECT id FROM user WHERE id = ? LIMIT 1"
)

type userHandler struct {
	db *sql.DB
}

// user create request length. This will limit how many data we read from
// request, to avoid attacks when someone might send large amounts of data.
const ucrLength = 100000

// userCreate() reads request, validates email, checks if user exists,
// saves user to db, and returns a JSON response.
func (uh *userHandler) userCreate(w http.ResponseWriter, r *http.Request) {
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

	user := &User{}
	// could not create user type from provided json
	if err = json.Unmarshal(body, user); err != nil {
		w.WriteHeader(422) // unprocessable entity
		APIResponse{Error: "Unprocessable entity"}.response(w)

		log.Fatalln(err)
		return
	}

	if len(user.Email) <= 8 || len(user.Email) > 40 {
		APIResponse{Error: "Invalid email length"}.response(w)
		return
	}

	if len(user.Username) <= 2 || len(user.Username) > 30 {
		APIResponse{Error: "Invalid username length"}.response(w)
		return
	}

	// email validation
	if _, err = mail.ParseAddress(user.Email); err != nil {
		APIResponse{Error: "Invalid email address"}.response(w)
		return
	}

	// check if there's a user with that username/email already
	exist, err := uh.userExists(user)
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
	user.CreationDate = time.Now()
	err = uh.insert(user)
	if err != nil {
		log.Println(err)
		w.Header().Set(`Status`, string(http.StatusInternalServerError))
		APIResponse{Error: "Could not insert user"}.response(w)
		return
	}

	// user created!
	w.WriteHeader(http.StatusCreated)
	user.Password = []byte("") // hide user password from response
	APIResponse{Message: "User created successfully!", Metadata: user}.response(w)
}

// userExists() queries database to find out if user already exists based on username and email.
func (uh *userHandler) userExists(u *User) (bool, error) {
	var id int

	// Prepare statement for reading data
	err := uh.db.QueryRow(findUserByNameOrEmail, u.Username, u.Email).Scan(&id)
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
// This is idempotent and in case it's called multiple times it will not encrypt the password multiple times.
func (u *User) encryptPassword() {
	if u.passwordEncoded {
		return
	}
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}

	dk, err := scrypt.Key([]byte(u.Password), salt, 16384, 8, 1, 32)
	if err != nil {
		log.Fatalln(err)
	}

	u.Password = dk
	u.Salt = salt
	u.passwordEncoded = true
}

// insert() saves newly created user in database
func (uh *userHandler) insert(u *User) error {
	stmt, err := uh.db.Prepare(insertUserQuery)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(u.Username, u.Password, u.Salt, u.Email, u.CreationDate)
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

// Checks if user exists by given id. Returns true if user exists, or false if not.
// Also returns error or nil if no error occurred.
func (uh *userHandler) findByID(id int) (bool, error) {
	err := uh.db.QueryRow(findUserByID, id).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			// something went wrong with this query
			return false, err
		}
		// user doesn't exist
		return false, nil
	}

	// user already exists
	return true, nil
}
