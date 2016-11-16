package routesHandlers

import (
	"controllers"
	"db"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
	"reflect"
)

//TaskHandler handles /tasks route and perform all CRUD operations based on method type
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getTaskHandler(w, r)
		return
	}
	if r.Method == http.MethodPost {
		insertTaskHandler(w, r)
		return
	}

	if r.Method == http.MethodPut {
		updateTaskHandler(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		removeTaskHandler(w, r)
		return
	}

	// creating an error json object to be passed to the http response
	newError := errorObj{Message: "Only Get,Delete,Put,Post requests are allowed", Resource: "Task"}
	json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Request method not allowed ", Status: http.StatusMethodNotAllowed}
	writeJSON(w, json)
	fmt.Println("Only Get,Delete,Put,Post requests are allowed", http.StatusMethodNotAllowed)

}

func insertTaskHandler(w http.ResponseWriter, r *http.Request) {
	newTaskData := models.Task{}
	// pass the memory address of the body object
	// this will populate the struct with the values from the request body
	// any field that is not in the request body will have its default value ex: for string it will be "" for arrays it will be []
	err := json.NewDecoder(r.Body).Decode(&newTaskData)
	if err != nil {
		newError := errorObj{Message: "Unable to parse request body . " + err.Error(), Resource: "Tasks"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Unable to parse request body . "+err.Error(), http.StatusBadRequest)
		return
	}
	// creating a taskcontroller and path to it a new db session
	taskController := controllers.NewTaskController(db.GetSession())
	// close db session after finishing working
	defer taskController.Session.Close()

	id, err := taskController.InsertTask(newTaskData)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Tasks"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}
	json := successJSONObj{Status: http.StatusCreated, Message: "Task added successfully", Results: map[string]string{"id": id.String()}}
	writeJSON(w, json)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// getting taskID from url passed parameters
	taskID := r.URL.Query().Get("id")
	// creating error json object to be passed with the response if the taskID is not provided
	if taskID == "" {
		newError := errorObj{Message: "ID of Task must be provided as a query parameter with key = id ex:(?id=taskID)", Resource: "Task"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of Task must be provided as a query parameter with key = id ex:(?id=taskID)", http.StatusBadRequest)
		return
	}

	// creating a taskcontroller and path to it a new db session
	taskController := controllers.NewTaskController(db.GetSession())
	// close db session after finishing working
	defer taskController.Session.Close()

	// calling taskController getTask and pass to the taksID
	task, err := taskController.GetTask(taskID)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Tasks"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successTaskJSONObj{Status: http.StatusOK, Message: "OK", Results: task}
	writeJSON(w, json)

}

func removeTaskHandler(w http.ResponseWriter, r *http.Request) {
	// getting taskID from url passed parameters
	taskID := r.URL.Query().Get("id")
	// creating error json object to be passed with the response if the taskID is not provided
	if taskID == "" {
		newError := errorObj{Message: "ID of Task must be provided as a query parameter with key = id ex:(?id=taskID)", Resource: "Task"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of Task must be provided as a query parameter with key = id ex:(?id=taskID)", http.StatusBadRequest)
		return
	}

	// creating a taskcontroller and path to it a new db session
	taskController := controllers.NewTaskController(db.GetSession())
	// close db session after finishing working
	defer taskController.Session.Close()

	err := taskController.RemoveTask(taskID)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Tasks"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successJSONObj{Status: http.StatusNoContent, Message: "Task deleted successfully"}
	writeJSON(w, json)

}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// getting taskID from url passed parameters
	taskID := r.URL.Query().Get("id")
	// creating error json object to be passed with the response if the taskID is not provided
	if taskID == "" {
		newError := errorObj{Message: "ID of Task must be provided as a query parameter with key = id ex:(?id=taskID)", Resource: "Task"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of Task must be provided as a query parameter with key = id ex:(?id=taskID)", http.StatusBadRequest)
		return
	}

	// creating a taskcontroller and path to it a new db session
	taskController := controllers.NewTaskController(db.GetSession())
	// close db session after finishing working
	defer taskController.Session.Close()

	updatedTaskdata := models.Task{}

	// pass the memory address of the body object
	// this will populate the struct with the values from the request body
	// any field that is not in the request body will have its default value ex: for string it will be "" for arrays it will be []
	err := json.NewDecoder(r.Body).Decode(&updatedTaskdata)

	if err != nil {
		newError := errorObj{Message: "Unable to parse request body . " + err.Error(), Resource: "Tasks"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Unable to parse request body . "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedTaskDataMap := make(map[string]string)
	// looping on struct and get all updated key value pairs passed by the user
	values := reflect.ValueOf(updatedTaskdata)
	fields := reflect.TypeOf(updatedTaskdata)
	tempintslice := []int{0}
	ielements := reflect.ValueOf(updatedTaskdata).NumField()

	for i := 0; i < ielements; i++ {
		v := values.Field(i).Interface()
		// fmt.Println(v)
		if v != "" {
			tempintslice[0] = i
			f := fields.FieldByIndex(tempintslice)
			// check if v is instanse of string
			_, ok := v.(string)
			if ok {
				updatedTaskDataMap[f.Name] = v.(string)
			}
		}
	}

	err = taskController.UpdateTask(updatedTaskDataMap, taskID)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Tasks"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successJSONObj{Status: http.StatusAccepted, Message: "Task is updated successfully"}
	writeJSON(w, json)

}
