package db

import mgo "gopkg.in/mgo.v2"

//GetSession connects to database and return a session
func GetSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}
