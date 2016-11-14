package models

import "gopkg.in/mgo.v2/bson"

type (
	// User represents the structure of our resource
	User struct {
		ID         bson.ObjectId `json:"id" bson:"_id"`
		FirstName  string        `json:"firstName" bson:"firstName"`
		LastName   string        `json:"lastName" bson:"lastName"`
		Email      string        `json:"email" bson:"email"`
		CalendarID string        `json:"calendarID" bson:"calendarID"`
	}
)
