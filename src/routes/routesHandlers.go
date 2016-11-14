package routes

import (
	"chatbot"
	"controllers"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// handles /welcome and respond with generated UUID
func handleWelcome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// creating an error json object to be passed to the http response
		newError := errorObj{Message: "Only Get requests are allowed", Resource: "Welcome Chat"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Request method not allowed ", Status: http.StatusMethodNotAllowed}
		writeJSON(w, json)
		fmt.Println("Only Get requests are allowed.", http.StatusMethodNotAllowed)
		return
	}
	res := chatbot.Welcome()
	json := successSingleJSONObj{Status: http.StatusOK, Message: "OK", Results: []map[string]string{res}}

	writeJSON(w, json)

}

// handle Handles / and respond with HTML Page
func handle(w http.ResponseWriter, r *http.Request) {
	body :=
		"<!DOCTYPE html><html><head><title>Chatbot</title></head><body><pre style=\"font-family: monospace;\">\n" +
			"Available Routes:\n\n" +
			"  GET  /welcome -> handleWelcome\n" +
			"  POST /chat    -> handleChat\n" +
			"  GET  /        -> handle        (current)\n" +
			"</pre></body></html>"
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintln(w, body)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	// Make sure only POST requests are handled
	if r.Method != http.MethodPost {
		newError := errorObj{Message: "Only POST requests are allowed", Resource: "Event Chat"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Request method not allowed ", Status: http.StatusMethodNotAllowed}
		writeJSON(w, json)
		fmt.Println("Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Make sure a UUID exists in the Authorization header
	uuid := r.Header.Get("Authorization")
	if uuid == "" {
		newError := errorObj{Message: "Missing or empty Authorization header.", Resource: "Event Chat"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "unAuthorized access", Status: http.StatusUnauthorized}
		writeJSON(w, json)
		fmt.Println("Missing or empty Authorization header.", http.StatusUnauthorized)
		return
	}

	isAuthenticated := chatbot.CheckIfAuthenticated(uuid)
	if !isAuthenticated {
		newError := errorObj{Message: "No session found for: " + uuid + " .", Resource: "Event Chat"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "unAuthorized access", Status: http.StatusUnauthorized}
		writeJSON(w, json)
		fmt.Println("No session found for: "+uuid+" .", http.StatusUnauthorized)
		return
	}

	// Parse the JSON string in the body of the request
	data := make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		newError := errorObj{Message: "Couldn't decode JSON: " + err.Error() + " .", Resource: "Event Chat"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Couldn't decode JSON: "+err.Error()+" .", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Make sure a message key is defined in the body of the request
	_, messageFound := data["message"]
	if !messageFound {
		newError := errorObj{Message: "Missing message key in body.", Resource: "Event Chat"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request ", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Missing message key in body.", http.StatusBadRequest)
		return
	}

	res, err := chatbot.Chat(uuid, data)

	if err != nil {
		newError := errorObj{Message: err.Error() + " .Unable to process entity . please try again", Resource: "Event Chat"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "StatusUnprocessableEntity", Status: 422}
		writeJSON(w, json)
		fmt.Println("Unable to process entity . please try again", 422)
		return
	}

	json := successSingleJSONObj{Status: http.StatusOK, Message: "OK", Results: []map[string]string{res}}
	writeJSON(w, json)

}

func eventListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// creating an error json object to be passed to the http response
		newError := errorObj{Message: "Only Get requests are allowed", Resource: "Google Calendar list events"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Request method not allowed ", Status: http.StatusMethodNotAllowed}
		writeJSON(w, json)
		fmt.Println("Only Get requests are allowed.", http.StatusMethodNotAllowed)
		return
	}
	// catching any error thrown
	defer func() {
		err := recover()
		if err != nil {
			newError := errorObj{Message: err.(string) + " . Unable to connect to Google calendar", Resource: "Google calendar"}
			json := errorsJSONObj{Errors: []errorObj{newError}, Message: "StatusUnprocessableEntity", Status: http.StatusInternalServerError}
			writeJSON(w, json)
			fmt.Println(err.(string)+" . Unable to connect to Google calendar", http.StatusInternalServerError)
		}
	}()

	calendarID := "k352nehms8mbf0hbe69jat2qig@group.calendar.google.com"
	eventController := controllers.NewEventController()
	calendarTitle, events, err := eventController.ListEvents(calendarID)

	if err != nil {
		// creating error json object to be send with the response
		newError := errorObj{Message: err.Error(), Resource: "Google calendar List Events"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successListJSONObj{Status: http.StatusOK, Message: "OK", Page: 1, CalendarTitle: calendarTitle, Results: events}
	writeJSON(w, json)

}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getEvent(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		deleteEvent(w, r)
		return
	}

	if r.Method == http.MethodPut {
		updateEvent(w, r)
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	// catching any error thrown
	defer func() {
		err := recover()
		if err != nil {
			newError := errorObj{Message: err.(string) + " . Unable to connect to Google calendar", Resource: "Google calendar"}
			json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
			writeJSON(w, json)
			fmt.Println(err.(string)+" . Unable to connect to Google calendar", http.StatusInternalServerError)
		}
	}()

	calendarID := "k352nehms8mbf0hbe69jat2qig@group.calendar.google.com"
	eventID := r.URL.Query().Get("id")

	// creating error json object to be passed with the response if the eventid is not provided
	if eventID == "" {
		newError := errorObj{Message: "ID of event must be provided as a query parameter with key = id ex:(?id=eventID)", Resource: "Google calendar Event"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of event must be provided as a query parameter with key = id ex:(?id=eventID)", http.StatusBadRequest)
		return
	}
	eventController := controllers.NewEventController()
	err := eventController.DeleteEvent(calendarID, eventID)
	if err != nil {
		newError := errorObj{Message: "Unable to delete event. " + err.Error(), Resource: "Google calendar Event"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println("Unable to delete event. "+err.Error(), http.StatusInternalServerError)
		return
	}

	json := successJSONObj{Status: http.StatusNoContent, Message: "Event deleted successfully"}
	writeJSON(w, json)

}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	// catching any error thrown
	defer func() {
		err := recover()
		if err != nil {
			newError := errorObj{Message: err.(string) + " . Unable to connect to Google calendar", Resource: "Google calendar"}
			json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
			writeJSON(w, json)
			fmt.Println(err.(string)+" . Unable to connect to Google calendar", http.StatusInternalServerError)
		}
	}()

	calendarID := "k352nehms8mbf0hbe69jat2qig@group.calendar.google.com"
	eventID := r.URL.Query().Get("id")

	// creating error json object to be passed with the response if the eventid is not provided
	if eventID == "" {
		newError := errorObj{Message: "ID of event must be provided as a query parameter with key = id ex:(?id=eventID)", Resource: "Google calendar Event"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of event must be provided as a query parameter with key = id ex:(?id=eventID)", http.StatusBadRequest)
		return
	}

	type updateEvent struct {
		Title            string   `json:"title"`
		Description      string   `json:"description"`
		StartDateTime    string   `json:"startDateTime"`
		EndDateTime      string   `json:"endDateTime"`
		Location         string   `json:"location"`
		OrganizerEmail   string   `json:"organizerEmail"`
		DeletedAttendees []string `json:"deletedAttendees"`
		NewAttendees     []string `json:"newAttendees"`
	}

	body := updateEvent{}

	// pass the memory address of the body object
	// this will populate the struct with the values from the request body
	// any field that is not in the request body will have its default value ex: for string it will be "" for arrays it will be []
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		newError := errorObj{Message: "Unable to parse request body . " + err.Error(), Resource: "Google calendar Event"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("Unable to parse request body . "+err.Error(), http.StatusBadRequest)
		return
	}
	eventMap := make(map[string]string)
	newAttendees := []string{}
	deletedAttendees := make(map[string]string)
	values := reflect.ValueOf(body)
	fields := reflect.TypeOf(body)
	tempintslice := []int{0}
	ielements := reflect.ValueOf(body).NumField()
	for i := 0; i < ielements; i++ {
		v := values.Field(i).Interface()
		if v != "" {
			tempintslice[0] = i
			f := fields.FieldByIndex(tempintslice)
			// check if v is instanse of string
			_, ok := v.(string)
			if ok {
				eventMap[f.Name] = v.(string)
			} else {
				// v is an arrays
				// check if the body field is NewAttendees and that the input array is not empty

				if f.Name == "NewAttendees" && len(v.([]string)) > 0 {
					// copy elements from the input array to newAttendees array . cannot use copy because the two arrays must have the same length
					// which is not applicable here

					newAttendees = v.([]string)
				}
				// check if the body field is deletedAttendees and that the input array is not empty
				fmt.Println(f.Name)

				if f.Name == "DeletedAttendees" && len(v.([]string)) > 0 {
					// iterate on the array and populate the deletedAttendees map
					for _, v := range v.([]string) {
						deletedAttendees[v] = v
					}
				}

			}
		}

	}
	fmt.Println(len(newAttendees), " ", len(deletedAttendees))
	eventController := controllers.NewEventController()
	updatedEvent, err := eventController.UpdateEvent(calendarID, eventID, newAttendees, deletedAttendees, eventMap)
	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Google calendar Event"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successEventJSONObj{Status: http.StatusOK, Message: "OK", Results: updatedEvent}
	writeJSON(w, json)

}

func getEvent(w http.ResponseWriter, r *http.Request) {
	// catching any error thrown
	defer func() {
		err := recover()
		if err != nil {
			newError := errorObj{Message: err.(string) + " . Unable to connect to Google calendar", Resource: "Google calendar"}
			json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
			writeJSON(w, json)
			fmt.Println(err.(string)+" . Unable to connect to Google calendar", http.StatusInternalServerError)
		}
	}()

	calendarID := "k352nehms8mbf0hbe69jat2qig@group.calendar.google.com"
	eventID := r.URL.Query().Get("id")

	// creating error json object to be passed with the response if the eventid is not provided
	if eventID == "" {
		newError := errorObj{Message: "ID of event must be provided as a query parameter with key = id ex:(?id=eventID)", Resource: "Google calendar Event"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Bad Request", Status: http.StatusBadRequest}
		writeJSON(w, json)
		fmt.Println("ID of event must be provided as a query parameter with key = id ex:(?id=eventID)", http.StatusBadRequest)
		return
	}
	// calling controller's get event function the return the event from google calendar api
	eventController := controllers.NewEventController()
	event, err := eventController.GetEvent(calendarID, eventID)
	if err != nil {
		newError := errorObj{Message: err.Error(), Resource: "Google calendar Event"}
		json := errorsJSONObj{Errors: []errorObj{newError}, Message: "Internal Server Error", Status: http.StatusInternalServerError}
		writeJSON(w, json)
		fmt.Println(err.Error(), http.StatusInternalServerError)
		return
	}

	json := successEventJSONObj{Status: http.StatusOK, Message: "OK", Results: event}
	writeJSON(w, json)
}
