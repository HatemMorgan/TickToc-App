package routesHandlers

import (
	"calendarAuth"
	"controllers"
	"db"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
	"reflect"
)

//UsersHandler handles /user route and perform all the crud operation based on http method
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getUserHandler(w, r)
		return
	}
	if r.Method == http.MethodPost {
		insertUserHandler(w, r)
		return
	}

	if r.Method == http.MethodPut {
		updateUserHandler(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		removeUserHandler(w, r)
		return
	}

	// creating an error json object to be passed to the http response
	newError := errorObj{Message: "Only Get,Delete,Put,Post requests are allowed", Resource: "User"}
	json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Request method not allowed ", Status: http.StatusMethodNotAllowed}
	writeJSON(w, json)
	fmt.Println("Only Get,Delete,Put,Post requests are allowed", http.StatusMethodNotAllowed)
}

func insertUserHandler(w http.ResponseWriter, r *http.Request) {
	newUserData := models.User{}
	// pass the memory address of the body object
	// this will populate the struct with the values from the request body
	// any field that is not in the request body will have its default value ex: for string it will be "" for arrays it will be []
	//
	tokenCode := r.Header.Get("tokenCode")
	if tokenCode == "" {
		authURL, err := calendarAuth.GetAuthURLFromWeb()
		if err != nil {
			newError := errorObj{Message: err.Error(), Resource: "Users"}
			json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
			fmt.Println(err.Error(), http.StatusInternalServerError)
			writeJSON(w, json)
			return
		}

		json := authURLJSONObj{Message: "OK", Status: http.StatusOK, Results: map[string]string{"authURL": authURL}}
		writeJSON(w, json)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&newUserData)
	if err != nil {
		newError := errorObj{Message: "Unable to parse request body . " + err.Error(), Resource: "Users"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Unable to parse request body . "+err.Error(), http.StatusBadRequest)
		return
	}

	// creating a userController and path to it a new db session
	userController := controllers.NewUserController(db.GetSession())
	// close db session after finishing working
	defer userController.Session.Close()

	userID, err := userController.InsertUser(newUserData, tokenCode)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Users"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := userController.GetUser(userID)
	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Users"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}
	if user.CalendarID == "" {
		eventController := controllers.NewEventController()
		calendarID, err := eventController.CreateAdvancedLabCalendar(user.Token)
		if err != nil {
			newError := errorObj{Message: err.Error(), Resource: "Google Calendar"}
			json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
			writeJSON(w, json)
			fmt.Println(err.Error(), http.StatusInternalServerError)
			return
		}

		userController.UpdateUser(map[string]string{"CalenarID": calendarID}, userID)
	}
	json := successJSONObj{Status: http.StatusCreated, Message: "user added successfully", Results: map[string]string{"userID": userID}}
	writeJSON(w, json)
}

func removeUserHandler(w http.ResponseWriter, r *http.Request) {
	// getting userID from url passed parameters
	userID := r.URL.Query().Get("id")
	// creating error json object to be passed with the response if the userID is not provided
	if userID == "" {
		newError := errorObj{Message: "ID of User must be provided as a query parameter with key = id ex:(?id=userID)", Resource: "User"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of User must be provided as a query parameter with key = id ex:(?id=userID)", http.StatusBadRequest)
		return
	}

	// creating a userController and path to it a new db session
	userController := controllers.NewUserController(db.GetSession())
	// close db session after finishing working
	defer userController.Session.Close()

	err := userController.RemoveUser(userID)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "User"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successJSONObj{Status: http.StatusNoContent, Message: "User deleted successfully"}
	writeJSON(w, json)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// getting userID from url passed parameters
	userID := r.URL.Query().Get("id")
	// creating error json object to be passed with the response if the userID is not provided
	if userID == "" {
		newError := errorObj{Message: "ID of User must be provided as a query parameter with key = id ex:(?id=userID)", Resource: "User"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of User must be provided as a query parameter with key = id ex:(?id=userID)", http.StatusBadRequest)
		return
	}

	// creating a userController and path to it a new db session
	userController := controllers.NewUserController(db.GetSession())
	// close db session after finishing working
	defer userController.Session.Close()

	updatedUserdata := models.User{}

	// pass the memory address of the body object
	// this will populate the struct with the values from the request body
	// any field that is not in the request body will have its default value ex: for string it will be "" for arrays it will be []
	err := json.NewDecoder(r.Body).Decode(&updatedUserdata)

	if err != nil {
		newError := errorObj{Message: "Unable to parse request body . " + err.Error(), Resource: "User"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Unable to parse request body . "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedUserDataMap := make(map[string]string)
	// looping on struct and get all updated key value pairs passed by the user
	values := reflect.ValueOf(updatedUserdata)
	fields := reflect.TypeOf(updatedUserdata)
	tempintslice := []int{0}
	ielements := reflect.ValueOf(updatedUserdata).NumField()

	for i := 0; i < ielements; i++ {
		v := values.Field(i).Interface()
		if v != "" {
			tempintslice[0] = i
			f := fields.FieldByIndex(tempintslice)
			// check if v is instanse of string
			_, ok := v.(string)
			if ok {
				updatedUserDataMap[f.Name] = v.(string)
			}
		}
	}

	err = userController.UpdateUser(updatedUserDataMap, userID)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "User"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successJSONObj{Status: http.StatusAccepted, Message: "User is updated successfully"}
	writeJSON(w, json)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// getting userID from url passed parameters
	userID := r.URL.Query().Get("id")
	// creating error json object to be passed with the response if the userID is not provided
	if userID == "" {
		newError := errorObj{Message: "ID of User must be provided as a query parameter with key = id ex:(?id=userID)", Resource: "User"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of User must be provided as a query parameter with key = id ex:(?id=userID)", http.StatusBadRequest)
		return
	}

	// creating a userController and path to it a new db session
	userController := controllers.NewUserController(db.GetSession())
	// close db session after finishing working
	defer userController.Session.Close()

	// calling taskController getTask and pass to the taksID
	user, err := userController.GetUser(userID)

	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Users"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal server error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successUserJSONObj{Status: http.StatusOK, Message: "OK", Results: user}
	writeJSON(w, json)
}

func getUser(userID string) (models.User, error) {
	// creating a userController and path to it a new db session
	userController := controllers.NewUserController(db.GetSession())
	// close db session after finishing working
	defer userController.Session.Close()

	// calling taskController getTask and pass to the taksID
	user, err := userController.GetUser(userID)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
