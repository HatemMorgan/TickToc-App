package chatbot

import (
	"controllers"
	"crypto/md5"
	"db"
	"encoding/hex"
	"fmt"
	"models"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"strconv"
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

func defaultProcessor(session Session, message string) (string, error) {
	return "", fmt.Errorf("You must use either the task or event chats")
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

//EventChat is called by chat/event route handler to save message and return the new question
func EventChat(uuid, userID string, data map[string]string) (map[string]string, error) {

	// gets user session
	session, _ := sessions[uuid]

	// add user id to the session
	session["userID"] = userID
	// set processor to EventChatProcessor
	ProcessFunc(eventChatProcessor)

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

//TaskChat is called by chat/task route handler to save answer of question and return the new question
func TaskChat(uuid, userID string, data map[string]string) (map[string]string, error) {

	// gets user session
	session, _ := sessions[uuid]
	session["userID"] = userID

	// set processor to TaskChatProcessor
	ProcessFunc(taskchatbotProcess)

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

var x = -1
var attendeesEmails []string

//TaskchatbotProcess is used for handling chat of tasks questions and answers
func taskchatbotProcess(session Session, message string) (string, error) {
	taskController := controllers.NewTaskController(db.GetSession())
	if strings.EqualFold(message, "add") {
		session = make(map[string]string)
		x = 0
	}

	if strings.EqualFold(message, "again") && x == 6 {
		x = -1
		session = make(map[string]string)
		return fmt.Sprintf("%s", "If you want to add events, type 'add'!"), nil
	}

	if strings.EqualFold(message, "done") && x == 6 {
		x = -1
		startDateTime, _ := strconv.ParseInt(session["startDateTime"], 10, 64)
		endDateTime, _ := strconv.ParseInt(session["endDateTime"], 10, 64)

		userID := bson.ObjectIdHex(session["userID"])

		newTask := models.Task{
			Title:         session["title"],
			Description:   session["description"],
			StartDateTime: startDateTime,
			EndDateTime:   endDateTime,
			Location: models.Location{
				Latitude:  session["Latitude"],
				Longitude: session["Longitude"],
			},
			UserID: userID,
		}
		taskController.InsertTask(newTask)
		return fmt.Sprintf("%s", "If you want to add another events, type 'add'!"), nil

	}

	switch x {
	case 0:
		x = 1
		return fmt.Sprintf("%s", "Please enter the title of your task"), nil
	case 1:
		session["title"] = message
		x = 2
		return fmt.Sprintf("%s", "Please enter the description of your task"), nil
	case 2:
		session["description"] = message
		x = 3
		return fmt.Sprintf("%s", "Please enter the start dateTime of your task"), nil
	case 3:
		session["startDateTime"] = message
		x = 4
		return fmt.Sprintf("%s", "Please enter the end dateTime of your task"), nil
	case 4:
		session["endDateTime"] = message
		x = 5
		return fmt.Sprintf("%s", "Please enter the location of the event ex:(Longitude,Latitude)"), nil
	case 5:
		location := strings.Split(message, ",")
		if len(location) != 2 {
			return "", fmt.Errorf("%s", "Invalid location entry it should be in this format ''Longitude,Latitude'' ")
		}
		session["Longitude"] = location[0]
		session["Latitude"] = location[1]

		var task = "Title: " + session["title"] + " , Description: " + session["description"] + " ,Start DateTime: " + session["startDateTime"] + " , End DateTime: " + session["endDateTime"] + " ,Location: { Longitude: " + session["Longitude"] + " , Latitude: " + session["Latitude"] + " }"

		x = 6

		return fmt.Sprintf("So your new task is " + task + " . Either type done to add it or type again to re-add it ."), nil

	default:
		return "", fmt.Errorf("%s", "Invalid text!")

	}
}

//EventChatProcessor is used for handling chat of events questions and answers
func eventChatProcessor(session Session, message string) (string, error) {
	eventController := controllers.NewEventController()
	if strings.EqualFold(message, "add") {
		session = make(map[string]string)
		x = 0
	}

	if strings.EqualFold(message, "again") && x == 8 {
		x = -1
		return fmt.Sprintf("%s", "If you want to add events, type 'add'!"), nil
	}

	if strings.EqualFold(message, "done") && x == 8 {
		x = -1
		fmt.Println(session)
		event, err := eventController.InsertEvent(session, attendeesEmails)
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		return fmt.Sprintf("%v", event), nil

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
