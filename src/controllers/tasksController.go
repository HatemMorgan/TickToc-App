package controllers

import (
	"models"

	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

func (taskController TaskController) insertTask(newTask models.Task) (models.Task, error) {
	// add an ID
	newTask.ID = bson.NewObjectId()

	// Write the task to mongo
	err := taskController.session.DB("advanced_computer_lab").C("tasks").Insert(newTask)
	if err != nil {
		return models.Task{}, err
	}

	return newTask, nil

}

//GetTask retrieves an individual task resource
func (taskController TaskController) GetTask(id string) (models.Task, error) {
	// Verify id is ObjectId, otherwise return error
	if !bson.IsObjectIdHex(id) {
		return models.Task{}, fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectID := bson.ObjectIdHex(id)

	// get task from mongo
	task := models.Task{}
	err := taskController.session.DB("advanced_computer_lab").C("tasks").FindId(objectID).One(&task)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

//RemoveTask removes an existing task resource
func (taskController TaskController) RemoveTask(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectID := bson.ObjectIdHex(id)

	err := taskController.session.DB("advanced_computer_lab").C("tasks").RemoveId(objectID)
	if err != nil {
		return err
	}

	return nil
}
