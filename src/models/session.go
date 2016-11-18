package models

import "gopkg.in/mgo.v2/bson"
import "time"

type (
	//Session is the model of stored in mongodb and used to give user session for chatting
	Session struct {
		ID        bson.ObjectId `json:"id" bson:"_id"`
		UUID      string        `json:"UUID" bson:"UUID"`
		UserID    bson.ObjectId `json:"userID" bson:"userID"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
		ExpireAt  time.Time     `json:"expireAt" bson:"expireAt"`
	}
)
