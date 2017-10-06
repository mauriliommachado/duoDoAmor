package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
)

var fSession mgo.Session

func Start() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	fSession = *session
	fmt.Println("Sessão do banco criada")
}

func GetCollectionUsers() (c *mgo.Collection) {
	s := fSession.Copy()
	c = s.DB("duolol").C("users")
	return c
}
