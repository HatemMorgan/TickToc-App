package main

import (
	"time"

	mgo "gopkg.in/mgo.v2"

	// Autoload environment variables in .env

	"models"

	"controllers"

	"fmt"

	_ "github.com/joho/godotenv/autoload"
)

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

func main() {

	// // Use the PORT environment variable
	// port := os.Getenv("PORT")
	// // Default to 3000 if no PORT environment variable was defined
	// if port == "" {
	// 	port = "3000"
	// }
	// // Start the server
	// fmt.Printf("Listening on port %s...\n", port)
	// log.Fatalln(routes.Routing(":" + port))

	// Manually Testing for events

	// a := time.Now().UnixNano() / int64(time.Millisecond)
	// fmt.Printf("%d \n", a)

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

	// manual testing for tasks
	newTask := models.Task{
		Title:         "Embedded Quiz",
		Description:   "Sheet 6 and 7 Embedded C",
		StartDateTime: time.Now().UnixNano() / int64(time.Millisecond),
		EndDateTime:   time.Now().UnixNano() / int64(time.Millisecond),
		Location: models.Location{
			Latitude:  "0.002",
			Longitude: "-0.23324",
		},
	}

	taskController := controllers.NewTaskController(getSession())
	task, err := taskController.InsertTask(newTask)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(task)

}
