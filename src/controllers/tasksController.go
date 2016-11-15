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

//InsertTask is responsible to add new task to database
func (taskController TaskController) InsertTask(newTask models.Task) (models.Task, error) {
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
func (taskController TaskController) GetTask(id string) (bson.ObjectId, error) {
	// Verify id is ObjectId, otherwise return error
	if !bson.IsObjectIdHex(id) {
		return "", fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectID := bson.ObjectIdHex(id)

	// get task from mongo
	task := models.Task{}
	err := taskController.session.DB("advanced_computer_lab").C("tasks").FindId(objectID).One(&task)

	if err != nil {
		return "", err
	}

	return task.ID, nil
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

//UpdateTask update an exsisting task
func (taskController TaskController) UpdateTask(updatedMap map[string]string, id string) error {

	// Verify id is ObjectId, otherwise return error
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectID := bson.ObjectIdHex(id)

	// creating a model to add to it the updated key value pairs
	model := bson.M{}

	// iterating on the updated map to updated the old task
	for key, value := range updatedMap {
		model[key] = value
	}
	// updating the old task by the new values
	err := taskController.session.DB("advanced_computer_lab").C("tasks").UpdateId(objectID, bson.M{"$set": model})
	if err != nil {
		return err
	}

	return nil
}
