package main

import mgo "gopkg.in/mgo.v2"

type tokenHandler struct {
	session *mgo.Session
}
