package routes

import (
	"chatbot"
	"encoding/json"
	"fmt"
	"net/http"
)

// handles /welcome and respond with generated UUID
func handleWelcome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	json := chatbot.Welcome()
	writeJSON(w, json, "200 OK")

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
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Make sure a UUID exists in the Authorization header
	uuid := r.Header.Get("Authorization")
	if uuid == "" {
		http.Error(w, "Missing or empty Authorization header.", http.StatusUnauthorized)
		return
	}

	isAuthenticated := chatbot.CheckIfAuthenticated(uuid)
	if !isAuthenticated {
		http.Error(w, fmt.Sprintf("No session found for: %v.", uuid), http.StatusUnauthorized)
		return
	}

	// Parse the JSON string in the body of the request
	data := make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Couldn't decode JSON: %v.", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Make sure a message key is defined in the body of the request
	_, messageFound := data["message"]
	if !messageFound {
		http.Error(w, "Missing message key in body.", http.StatusBadRequest)
		return
	}

	json, err := chatbot.Chat(uuid, data)

	if err != nil {
		http.Error(w, err.Error(), 422 /* http.StatusUnprocessableEntity */)
		return
	}

	writeJSON(w, json, "200 OK")

}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// catching any error thrown
		defer func() {
			err := recover()
			if err != nil {
				http.Error(w, "Unable to connect to google calendar . Please try again. ", http.StatusInternalServerError)
			}
		}()

		// calendarID := "k352nehms8mbf0hbe69jat2qig@group.calendar.google.com"
		// eventID := r.URL.Query().Get("id")
		// calling controller's get event function the return the event from google calendar api
		// event, err := googleCalendarcontroller.GetEvent(calendarID, eventID)
		// if err != nil {
		// 	http.Error(w, "Unable to fetch event from google calendar . Please try again. ", http.StatusInternalServerError)
		// 	return
		// }

	}
}
