package controllers

import mgo "gopkg.in/mgo.v2"

type (

	// TaskController represents the controller for operating on the Task resource
	TaskController struct {
		session *mgo.Session
	}
)

//NewTaskController provides a reference to a TaskController with provided mongo session
func NewTaskController(s *mgo.Session) *TaskController {
	return &TaskController{s}
}
