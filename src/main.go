package main

import

// Autoload environment variables in .env

(
	"fmt"
	"log"
	"os"
	"routes"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "4000"
	}
	// Start the server
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatalln(routes.Routing(":" + port))

	//--------------------------------------------------------------------------------------------------------------------------------
	// hasher := md5.New()
	// hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	// uuid := hex.EncodeToString(hasher.Sum(nil))
	// userID := "582bf5e88a4e9e1d45dbdf05"
	// sessionModel := controllers.NewSessionModel(db.GetSession())
	// // sessionModel.InsertNewSession(uuid, userID)

	// session, err := sessionModel.GetSession("582bc3458a4e9e29e1a54439")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(session)
	// }
	// fmt.Println("--------------------------------------------------------------------")
	// session, err = sessionModel.GetSession("582bf5e88a4e9e1d45dbdf05")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(session)
	// }

	// fmt.Println(strconv.FormatInt(time.Now().Unix(), 10))
	// fmt.Println(time.Now())
	// fmt.Println(time.Now().UTC())

	// --------------------------------------------------------------------------------------------------------------------------------
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

	// manual testing for taskss

	// taskController := controllers.NewTaskController(db.GetSession())
	// tasks, err := taskController.ListTasks("582bc3458a4e9e29e1a54439")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(tasks)

	// now := time.Now()
	// newTask := models.Task{
	// 	Title:         "Project Advanced Computer Lab Milstone 2",
	// 	Description:   "Android Application",
	// 	StartDateTime: now.UTC(),
	// 	EndDateTime:   now.Add(120 * time.Hour).UTC(),
	// 	Location: models.Location{
	// 		Latitude:  "0.002",
	// 		Longitude: "-0.23324",
	// 	},
	// 	UserID: bson.ObjectIdHex("582bc3458a4e9e29e1a54439"),
	// }
	// id, err := taskController.InsertTask(newTask)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(id)

	// a := time.Now().UnixNano() / int64(time.Millisecond)
	// updateTaskMap := map[string]string{"Title": "Project Advanced Computer Lab deadline", "EndDateTime": strconv.FormatInt(a, 10), "Longitude": "0.222223"}
	// err := taskController.UpdateTask(updateTaskMap, "582bbb6b8a4e9e46c7df713e")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// task, err := taskController.GetTask("58412cb78a4e9e1115464f2a")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(task)

	// err := taskController.RemoveTask("582bb9878a4e9e30c301f184")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// manual test for users
	// userController := controllers.NewUserController(getSession())

	// newUser := models.User{
	// 	FirstName:  "Hatem",
	// 	LastName:   "Morgan",
	// 	Email:      "hatemmorgan17@gmail.com",
	// 	CalendarID: "k352nehms8mbf0hbe69jat2qig@group.calendar.google.com",
	// }
	// id, err := userController.InsertTask(newUser)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(id)

	// user, err := userController.GetUser("582bc2878a4e9e228731ad56")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(user)

	// err := userController.RemoveUser("582bc2878a4e9e228731ad56")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// updateUserMap := map[string]string{"lastName": "Elsayed", "ko": "ok"}
	// err := userController.UpdateUser(updateUserMap, "582bc3458a4e9e29e1a54439")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
