package models

import "gopkg.in/mgo.v2/bson"

type (
	//Task represents the structure of Task resource
	Task struct {
		ID            bson.ObjectId `json:"id" bson:"_id"`
		Title         string        `json:"title" bson:"title"`
		Description   string        `json:"description" bson:"description"`
		StartDateTime int64         `json:"startDateTime" bson:"startDateTime"`
		EndDateTime   int64         `json:"endDateTime" bson:"endDateTime"`
		Location      Location      `json:"location" bson:"location"`
	}

	//Location represent location of a task
	Location struct {
		Latitude  string `json:"latitude" bson:"latitude"`
		Longitude string `json:"longitude" bson:"longitude"`
	}
)
