package controllers

import (
	"fmt"
	"models"

	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	//SessionModel represent the session model that operate on session resource
	SessionModel struct {
		DBSession *mgo.Session
	}
)

// db.sessions.createIndex({"expireAt":1},{expireAfterSeconds:0}) to create an expiry index on expireAt field of sessions
//  db.sessions.createIndex({"userID":1}) to create an index on userID field of sessions

//NewSessionModel provides a reference to a SessionModel with provided mongo session
func NewSessionModel(s *mgo.Session) *SessionModel {
	return &SessionModel{s}
}

//InsertNewSession is responsible to add new session to database
func (sessionModel SessionModel) InsertNewSession(UUID string, userID string) (bson.ObjectId, error) {
	newSession := models.Session{}
	// add an ID
	newSession.ID = bson.NewObjectId()

	// adding expiry date to the session 1 hour
	now := time.Now()
	newSession.CreatedAt = now.UTC()

	expireAt := now.Add(1 * time.Hour)
	newSession.ExpireAt = expireAt.UTC()

	// adding UUID and userID to newSession
	newSession.UUID = UUID

	// Verify id is ObjectId, otherwise return error
	if !bson.IsObjectIdHex(userID) {
		return "", fmt.Errorf("Invalid user ID")
	}
	// Grab id
	objectUserID := bson.ObjectIdHex(userID)
	newSession.UserID = objectUserID

	// Write the new session to mongo
	err := sessionModel.DBSession.DB("advanced_computer_lab").C("sessions").Insert(newSession)
	if err != nil {
		return "", fmt.Errorf("Unable to add new Session . %v ", err)
	}

	return newSession.ID, nil

}

//GetSession retrieves an individual session resource
func (sessionModel SessionModel) GetSession(userID string) (models.Session, error) {
	// Verify id is ObjectId, otherwise return error
	if !bson.IsObjectIdHex(userID) {
		return models.Session{}, fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectUserID := bson.ObjectIdHex(userID)

	// get session from mongo
	session := models.Session{}
	err := sessionModel.DBSession.DB("advanced_computer_lab").C("sessions").Find(bson.M{"userID": objectUserID}).One(&session)

	if err != nil {
		return models.Session{}, fmt.Errorf("No session uuid available for this user: %s . %v", userID, err)
	}

	return session, nil
}
