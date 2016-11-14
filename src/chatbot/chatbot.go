package chatbot

import (
	"controllers"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// WelcomeMessage A constant to hold the welcome message
	WelcomeMessage = "Tick-tock, Whenever you want to add an event, just type 'add'!"

	// sessions = {
	//   "uuid1" = Session{
	//     "history" = [
	//       "Message 1",
	//       "Message 2",
	//       "Message 3",
	//       ...
	//     ]
	//   },
	//   ...
	// }
	sessions = map[string]Session{}

	processor = defaultProcessor
)

type (
	// Session Holds info about a session
	Session map[string]string

	// JSON Holds a JSON object
	// JSON map[string]interface{}

	// Processor Alias for Process func
	Processor func(session Session, message string) (string, error)
)

var x = -1
var attendeesEmails []string

func defaultProcessor(session Session, message string) (string, error) {

	if strings.EqualFold(message, "add") {
		session = make(map[string]string)
		x = 0
	}

	if strings.EqualFold(message, "again") {
		x = -1
		return fmt.Sprintf("%s", "If you want to add events, type 'add'!"), nil
	}

	if strings.EqualFold(message, "done") {
		x = -1
		controllers.InsertEvent(session, attendeesEmails)
		return fmt.Sprintf("%s", "If you want to add another events, type 'add'!"), nil

	}

	switch x {
	case 0:
		x = 1
		return fmt.Sprintf("%s", "Please enter the title of the event"), nil
	case 1:
		session["title"] = message
		x = 2
		return fmt.Sprintf("%s", "Please enter the description of the event"), nil
	case 2:
		session["description"] = message
		x = 3
		return fmt.Sprintf("%s", "Please enter the start dateTime of the event"), nil
	case 3:
		session["startDateTime"] = message
		x = 4
		return fmt.Sprintf("%s", "Please enter the end dateTime of the event"), nil
	case 4:
		session["endDateTime"] = message
		x = 5
		return fmt.Sprintf("%s", "Please enter the location of the event"), nil
	case 5:
		session["location"] = message
		x = 6
		return fmt.Sprintf("%s", "Please enter the organizer email of the event"), nil
	case 6:
		session["organizerEmail"] = message
		x = 7
		return fmt.Sprintf("%s", "Please enter the attendees email of the event and split the emails with - "), nil
	case 7:
		session["attendeesEmails"] = message
		x = 8
		return fmt.Sprintf("%s", "Please choose a calendar to add the event to it"), nil
	case 8:
		session["calenderID"] = message

		attendeesEmails = strings.Split(session["attendeesEmails"], "-")
		fmt.Println(len(attendeesEmails))
		var attendees = " "

		for i, v := range attendeesEmails {
			attendees += " " + strconv.Itoa(i+1) + "- " + v + " "
		}

		var event = "Title: " + session["title"] + " , Description: " + session["description"] + " ,Start DateTime: " + session["startDateTime"] + " , End DateTime: " + session["endDateTime"] + " ,Location: " + session["location"] + " , Organizer email: " + session["organizerEmail"] + " , Attendees emails: " + attendees + " , Calender type: " + session["calenderID"]

		return fmt.Sprintf("So your event is " + event + " . Either type done to add it or type again to re-add it ."), nil
		// return fmt.Sprintf("%s", "This event is done! Either type 'add' or 'done'!"), nil

	default:
		return "", fmt.Errorf("%s", "Invalid text!")

	}

	// 	if strings.EqualFold(message, "chatbot") {
	// 		return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
	// 	}

	// 	return fmt.Sprintf("Hello %s, my name is chatbot. What was yours again?", message), nil

}

// ProcessFunc Sets the processor of the chatbot
func ProcessFunc(p Processor) {
	processor = p
}

//Welcome is called by welcome handler route to generate an UUID
func Welcome() map[string]string {
	// Generate a UUID.
	// bygeeb time stamp unix format and represent it in base 10
	hasher := md5.New()
	hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	uuid := hex.EncodeToString(hasher.Sum(nil))

	// Create a session for this UUID
	sessions[uuid] = Session{}

	// Write a JSON containg the welcome message and the generated UUID
	return map[string]string{
		"uuid":    uuid,
		"message": WelcomeMessage,
	}
}

//CheckIfAuthenticated checks if the user has a session opened and his uuid is valid
func CheckIfAuthenticated(uuid string) bool {
	// Make sure a session exists for the extracted UUID
	_, sessionFound := sessions[uuid]
	return sessionFound
}

//Chat is called by chat route handler to save message and return the new question
func Chat(uuid string, data map[string]string) (map[string]string, error) {

	// gets user session
	session, _ := sessions[uuid]

	// Process the received message
	message, err := processor(session, data["message"])
	if err != nil {
		return nil, err
	}

	// Write a JSON containg the processed response
	return map[string]string{
		"message": message,
	}, nil

}
