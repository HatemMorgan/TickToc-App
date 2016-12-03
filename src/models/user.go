package models

import "gopkg.in/mgo.v2/bson"
import "golang.org/x/oauth2"

type (
	// User represents the structure of our resource
	User struct {
		ID         bson.ObjectId `json:"id" bson:"_id"`
		FirstName  string        `json:"firstName" bson:"firstName"`
		LastName   string        `json:"lastName" bson:"lastName"`
		Email      string        `json:"email" bson:"email"`
		CalendarID string        `json:"calendarID" bson:"calendarID"`
		Token      *oauth2.Token `json:"token" bson:"token"`
	}
)

//IsValidField checks if fieldName is a valid field or not
func (user User) IsValidField(fieldName string) bool {
	fieldsMap := map[string]string{"firstName": "FirstName", "lastName": "LastName", "email": "Email", "calenarID": "CalendarID"}
	_, ok := fieldsMap[fieldName]
	return ok
}
