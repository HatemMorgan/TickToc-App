package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	//Task represents the structure of Task resource
	Task struct {
		ID            bson.ObjectId `json:"id" bson:"_id"`
		Title         string        `json:"title" bson:"title"`
		Description   string        `json:"description" bson:"description"`
		StartDateTime time.Time     `json:"startDateTime" bson:"startDateTime"`
		EndDateTime   time.Time     `json:"endDateTime" bson:"endDateTime"`
		Location      Location      `json:"location" bson:"location"`
	}

	//Location represent location of a task
	Location struct {
		Latitude  string `json:"latitude" bson:"latitude"`
		Longitude string `json:"longitude" bson:"longitude"`
	}
)
