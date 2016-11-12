package main

import (
	"chatbot"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Autoload environment variables in .env
import _ "github.com/joho/godotenv/autoload"

var x = -1

func chatbotProcess(session chatbot.Session, message string) (string, error) {

	var key string
	if strings.EqualFold(message, "add") {
		x = 0
		key = "event"
		session[key] = []string{}
	}

	if strings.EqualFold(message, "again") {
		x = -1
		return fmt.Sprintf("%s", "If you want to add events, type 'add'!"), nil
	}

	if strings.EqualFold(message, "done") {
		x = -1
		return fmt.Sprintf("%s", "If you want to add events, type 'add'!"), nil

	}

	switch x {
	case 0:
		x = 1
		return fmt.Sprintf("%s", "Please enter the title of the event"), nil
	case 1:
		session[key] = append(session[key], message)
		x = 2
		return fmt.Sprintf("%s", "Please enter the description of the event"), nil
	case 2:
		session[key] = append(session[key], message)
		x = 3
		return fmt.Sprintf("%s", "Please enter the start dateTime of the event"), nil
	case 3:
		session[key] = append(session[key], message)
		x = 4
		return fmt.Sprintf("%s", "Please enter the end dateTime of the event"), nil
	case 4:
		session[key] = append(session[key], message)
		x = 5
		return fmt.Sprintf("%s", "Please enter the location of the event"), nil
	case 5:
		session[key] = append(session[key], message)
		x = 6
		return fmt.Sprintf("%s", "Please enter the organizer email of the event"), nil
	case 6:
		session[key] = append(session[key], message)
		x = 7
		return fmt.Sprintf("%s", "Please enter the attendees email of the event and split the emails with - "), nil
	case 7:
		session[key] = append(session[key], message)
		x = 8
		return fmt.Sprintf("%s", "Please choose a calendar to add the event to it"), nil
	case 8:
		session[key] = append(session[key], message)
		var eventArray = session[key]
		var attendeesEmails = strings.Split(eventArray[6], "-")
		var attendees = " "

		for i, v := range attendeesEmails {
			attendees += " " + strconv.Itoa(i+1) + "- " + v + " "
		}

		var event = "Title: " + eventArray[0] + " , Description: " + eventArray[1] + " ,Start DateTime: " + eventArray[2] + " , End DateTime: " + eventArray[3] + " ,Location: " + eventArray[4] + " , Organizer email: " + eventArray[5] + " , Attendees emails: " + attendees + " , Calender type: " + eventArray[7]

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

func main() {
	// Uncomment the following lines to customize the chatbot
	//chatbot.WelcomeMessage = "What's your name?"
	chatbot.WelcomeMessage = "Tick-tock, Whenever you want to add an event, just type 'add'!"
	chatbot.ProcessFunc(chatbotProcess)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "3000"
	}

	// Start the server
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatalln(chatbot.Engage(":" + port))
}
