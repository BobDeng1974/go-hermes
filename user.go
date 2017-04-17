package main

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"time"

	mathRand "math/rand"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"golang.org/x/crypto/scrypt"
)

type userHandler struct {
	session *mgo.Session
}

// user create request length. This will limit how many data we read from
// request, to avoid attacks when someone might send large amounts of data.
const ucrLength = 100000

// userCreate() reads request, validates email, checks if user exists,
// saves user to db, and returns a JSON response.
func (uh *userHandler) userCreate(w http.ResponseWriter, r *http.Request) {
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
	user.generateRandomToken()

	// save user to db
	user.CreationDate = time.Now()
	user.Roles = []ACLRole{RoleUser}
	err = uh.insert(user)
	if err != nil {
		log.Println(err)
		w.Header().Set(`Status`, string(http.StatusInternalServerError))
		APIResponse{Error: "Could not insert user"}.response(w)
		return
	}

	// user created!
	w.WriteHeader(http.StatusCreated)
	user.Password = nil // hide user password from response
	APIResponse{Message: "User created successfully!", Metadata: user}.response(w)
}

// userExists() queries database to find out if user already exists based on username and email.
func (uh *userHandler) userExists(u *User) (bool, error) {
	result := User{}

	c := uh.session.DB("app").C("user")
	err := c.Find(bson.M{"$or": []bson.M{bson.M{"username": u.Username}, bson.M{"email": u.Email}}}).One(&result)
	if err == mgo.ErrNotFound {
		// user doesn't exist
		return false, nil
	}

	if err != nil {
		return false, err
	}

	// user already exists
	return true, nil
}

// random generator found on http://stackoverflow.com/a/22892986/1294631
func (u *User) generateRandomToken() {
	if u.tokenGenerated {
		return
	}

	n := 64
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}

	u.Token = string(b)
	u.tokenGenerated = true
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
	c := uh.session.DB("app").C("user")
	u.ID = bson.NewObjectId()

	err := c.Insert(bson.M{
		"_id":          u.ID,
		"username":     u.Username,
		"password":     u.Password,
		"salt":         u.Salt,
		"email":        u.Email,
		"roles":        u.Roles,
		"token":        u.Token,
		"creationDate": u.CreationDate,
	})

	if err != nil {
		return err
	}

	return nil
}

// Checks if user exists by given id. Returns true if user exists, or false if not.
// Also returns error or nil if no error occurred.
func (uh *userHandler) findByID(id bson.ObjectId) (bool, error) {
	var u User
	c := uh.session.DB("app").C("user")
	err := c.FindId(id).One(&u)
	// user not found
	if err == mgo.ErrNotFound {
		return false, nil
	}

	// something went wrong
	if err != nil {
		return false, err
	}

	// user already exists
	return true, nil
}
