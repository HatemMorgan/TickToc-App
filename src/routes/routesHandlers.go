package routes

import (
	"chatbot"
	"encoding/json"
	"fmt"
	"googleCalendarcontroller"
	"net/http"
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

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
		eventID := "gc23i3fr8kq9nph2c9nknbvvtc" //r.URL.Query().Get("id")
		// calling controller's get event function the return the event from google calendar api
		event, err := googleCalendarcontroller.GetEvent(calendarID, eventID)
		if err != nil {
			newError := errorObj{Message: err.Error(), Resource: "Google calendar Event"}
			json := errorsJSONObj{Errors: []errorObj{newError}, Message: "StatusUnprocessableEntity", Status: http.StatusInternalServerError}
			writeJSON(w, json)
			fmt.Println(err.Error(), http.StatusInternalServerError)
			return
		}

		json := successEventJSONObj{Status: http.StatusOK, Message: "OK", Results: event}
		writeJSON(w, json)

	}
}
