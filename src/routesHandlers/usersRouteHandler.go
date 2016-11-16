package routesHandlers

import (
	"controllers"
	"db"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {

}

func insertUserHandler(w http.ResponseWriter, r *http.Request) {
	newUserData := models.User{}
	// pass the memory address of the body object
	// this will populate the struct with the values from the request body
	// any field that is not in the request body will have its default value ex: for string it will be "" for arrays it will be []
	err := json.NewDecoder(r.Body).Decode(&newUserData)
	if err != nil {
		newError := errorObj{Message: "Unable to parse request body . " + err.Error(), Resource: "Users"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Unable to parse request body . "+err.Error(), http.StatusBadRequest)
		return
	}
	// creating a taskcontroller and path to it a new db session
	userController := controllers.NewUserController(db.GetSession())
	// close db session after finishing working
	defer userController.Session.Close()

	id, err := taskController.InsertTask(newTaskData)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Tasks"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}
	json := successJSONObj{Status: http.StatusCreated, Message: "Taks added successfully", Results: map[string]string{"id": id.String()}}
	writeJSON(w, json)
}

func removeUserHandler(w http.ResponseWriter, r *http.Request) {

}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func getUserHandler(w http.ResponseWriter, r *http.Request) {

}
