package controllers

import (
	"models"

	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (

	//UserController represents the controller for operating on the user resource
	UserController struct {
		session *mgo.Session
	}
)

//NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

//InsertTask is responsible to add new task to database
func (userController UserController) InsertTask(newUser models.User) (bson.ObjectId, error) {
	// add an ID
	newUser.ID = bson.NewObjectId()

	// Write the user to mongo
	err := userController.session.DB("advanced_computer_lab").C("users").Insert(newUser)
	if err != nil {
		return "", fmt.Errorf("Unable to add new Task . %v ", err)
	}

	return newUser.ID, nil
}

//GetUser retrieves an individual user resource
func (userController UserController) GetUser(id string) (models.User, error) {
	// Verify id is ObjectId, otherwise return error
	if !bson.IsObjectIdHex(id) {
		return models.User{}, fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectID := bson.ObjectIdHex(id)

	// get user from mongo
	user := models.User{}
	err := userController.session.DB("advanced_computer_lab").C("users").FindId(objectID).One(&user)

	if err != nil {
		return models.User{}, fmt.Errorf("Unable to get user with id: %s . %v", id, err)
	}

	return user, nil
}

//RemoveUser removes an existing user resource
func (userController UserController) RemoveUser(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("Invalid ID")
	}
	// Grab id
	objectID := bson.ObjectIdHex(id)

	err := userController.session.DB("advanced_computer_lab").C("users").RemoveId(objectID)
	if err != nil {
		return fmt.Errorf("Unable to remove user with id: %s . %v", id, err)
	}

	return nil
}

//UpdateUser update an exsisting user
func (userController UserController) UpdateUser(updatedMap map[string]string, id string) error {

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
		// make sure that the field is a valid field for user resource
		fieldsMap := map[string]string{"firstName": "FirstName", "lastName": "LastName", "email": "Email", "calenarID": "CalendarID"}
		_, ok := fieldsMap[key]
		if !ok {
			return fmt.Errorf("Invalid Field with this name: %s", key)
		}

		// adding key value pairs
		model[key] = value
	}

	// updating the old user by the new values
	err := userController.session.DB("advanced_computer_lab").C("users").UpdateId(objectID, bson.M{"$set": model})
	if err != nil {
		return err
	}

	return nil
}
