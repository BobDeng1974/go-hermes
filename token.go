package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type tokenHandler struct {
	session *mgo.Session
}

func (th *tokenHandler) getUserByToken(token string) (User, error) {
	var user User
	c := th.session.DB("app").C("user")
	err := c.Find(bson.M{"token": token}).One(&user)
	// no user found by this token, or something went wrong
	if err != nil {
		return user, err
	}

	return user, nil
}
