package main

import (
	"GoogleCalendarcontroller"
	"chatbot"
	"fmt"
	"log"
	"os"
	"routes"
	"strconv"
	"strings"

	// Autoload environment variables in .env

	_ "github.com/joho/godotenv/autoload"
)

var x = -1
var attendeesEmails []string

func chatbotProcess(session chatbot.Session, message string) (string, error) {

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
		GoogleCalendarcontroller.InsertEvent(session, attendeesEmails)
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

func main() {
	// Uncomment the following lines to customize the chatbot
	// chatbot.WelcomeMessage = "What's your name?"
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
	log.Fatalln(routes.Routing(":" + port))

	// testMap := make(map[string]string)
	// testMap["title"] = "Cairo Party"
	// testMap["description"] = "lets party"
	// testMap["startDateTime"] = "2016-11-13T22:00:00-07:00"
	// testMap["endDateTime"] = "2016-11-13T23:00:00-07:00"
	// testMap["location"] = "Cairo,Egypt"
	// testMap["organizerEmail"] = "hatemmorgan17@gmail.com"
	// testMap["calenderID"] = "k352nehms8mbf0hbe69jat2qig@group.calendar.google.com"
	// attendees := []string{"hatemmorgan17@gmail.com", "omartagguv@gmail.com", "jojo@gmail.com"}
	// GoogleCalendarcontroller.InsertEvent(testMap, attendees) //event id = gc23i3fr8kq9nph2c9nknbvvtc

	//GoogleCalendarcontroller.CreateAdvancedLabCalendar()  // created calendar id = k352nehms8mbf0hbe69jat2qig@group.calendar.google.com

	// GoogleCalendarcontroller.DeleteCalendar("b3cjs1oc7ql5jaecm98hv3fke0@group.calendar.google.com")

	// GoogleCalendarcontroller.ListEvents("k352nehms8mbf0hbe69jat2qig@group.calendar.google.com")

	// GoogleCalendarcontroller.GetEvent("k352nehms8mbf0hbe69jat2qig@group.calendar.google.com", "qobjl5rj6ebi2vuhukbli5oamk")

	// updatedEvent := make(map[string]string)
	// updatedEvent["title"] = "Dream Park"
	// newAttendees := []string{"Ahmed@gmail.com"}
	// deletedAttendees := make(map[string]string)
	// deletedAttendees["omartagguv@gmail.com"] = "omartagguv@gmail.com"
	// GoogleCalendarcontroller.UpdateEvent("k352nehms8mbf0hbe69jat2qig@group.calendar.google.com", "j518p4bcagq8kt1717vvmb8bf0", newAttendees, deletedAttendees, updatedEvent)
	// controller.GetControllerList()

	// GoogleCalendarcontroller.DeleteEvent("k352nehms8mbf0hbe69jat2qig@group.calendar.google.com", "j518p4bcagq8kt1717vvmb8bf0")
}
