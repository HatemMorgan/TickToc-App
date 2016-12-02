package controllers

import (
	"models"
	"time"

	"fmt"

	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (

	// TaskController represents the controller for operating on the Task resource
	TaskController struct {
		Session *mgo.Session
	}
)

//NewTaskController provides a reference to a TaskController with provided mongo session
func NewTaskController(s *mgo.Session) *TaskController {
	return &TaskController{s}
}

//InsertTask is responsible to add new task to database
func (taskController TaskController) InsertTask(newTask models.Task) (bson.ObjectId, error) {
	// add an ID
	newTask.ID = bson.NewObjectId()
	// Write the task to mongo
	err := taskController.Session.DB("advanced_computer_lab").C("tasks").Insert(newTask)
	if err != nil {
		return "", fmt.Errorf("Unable to add new Task . %v ", err)
	}

	return newTask.ID, nil

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
	err := taskController.Session.DB("advanced_computer_lab").C("tasks").FindId(objectID).One(&task)

	if err != nil {
		return models.Task{}, fmt.Errorf("Unable to get task with id: %s . %v", id, err)
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

	err := taskController.Session.DB("advanced_computer_lab").C("tasks").RemoveId(objectID)
	if err != nil {
		return fmt.Errorf("Unable to remove task with id: %s . %v", id, err)
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

		// make sure that the field is a valid field for Task resource
		fieldsMap := map[string]string{"Title": "title", "Description": "description", "StartDateTime": "startDateTime", "EndDateTime": "endDateTime", "Latitude": "latitude", "Longitude": "longitude"}
		if _, ok := fieldsMap[key]; !ok {
			return fmt.Errorf("Invalid Field with this name: %s", key)
		}
		// check if the updated value is the longitude of location to update the location object
		if key == "Longitude" {
			model["location.longitude"] = value
			continue
		}

		if key == "Latitude" {
			model["location.latitude"] = value
			continue
		}

		// casting from string to int64 because the data type of startdatetime and enddatetime is int64
		if key == "StartDateTime" || key == "EndDateTime" {
			num, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("Wrong Date format . Date must be converted to milliseconds (long int) ")
			}
			if num != 0 {
				model[fieldsMap[key]] = num
			}
			continue

		}
		// if the key is not longitude or latitude so it is a field in the document being updated so update the value
		// of the field crossponding to the key given
		k := fieldsMap[key]
		model[k] = value
	}
	// fmt.Println(model)
	// updating the old task by the new values
	err := taskController.Session.DB("advanced_computer_lab").C("tasks").UpdateId(objectID, bson.M{"$set": model})
	if err != nil {
		return err
	}

	return nil
}

//ListTasks lists all tasks
func (taskController TaskController) ListTasks(id string) ([]models.Task, error) {
	// Verify id is ObjectId, otherwise return error
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectID := bson.ObjectIdHex(id)

	// get tasks from mongo

	now := time.Now().UTC()

	tasks := []models.Task{}
	err := taskController.Session.DB("advanced_computer_lab").C("tasks").Find(bson.M{"_id": objectID, "endDateTime": bson.M{"$gte": now}}).All(&tasks)

	if err != nil {
		return nil, fmt.Errorf("Unable to get task with id: %s . %v", id, err)
	}

	return tasks, nil
}
