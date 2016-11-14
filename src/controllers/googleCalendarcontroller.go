package controllers

import (
	"calendarAuth"
	"log"

	"fmt"

	"strings"

	calendar "google.golang.org/api/calendar/v3"
)

//GetCallenderList returns a list of all Calendars
func GetCallenderList() {
	srv, err := calendarAuth.GetCalendarService()
	listRes, err := srv.CalendarList.List().Fields("items/id").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of calendars: %v", err)
	}
	for _, v := range listRes.Items {
		log.Printf("Calendar ID: %v and description: \n", v.Id)
	}

	log.Println("------------------------------------------ ")

	if len(listRes.Items) > 0 {
		id := listRes.Items[2].Id
		res, err := srv.Events.List(id).Fields("items(updated,summary)", "summary", "nextPageToken").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve calendar events list: %v", err)
		}
		for _, v := range res.Items {
			log.Printf("Calendar ID %q event: %v: %q\n", id, v.Updated, v.Summary)
		}
		log.Printf("Calendar ID %q Summary: %v\n", id, res.Summary)
		log.Printf("Calendar ID %q next page token: %v\n", id, res.NextPageToken)
	}
}

// InsertEvent takes the new event entry from user and create a new event then insert it using google calendar api
func InsertEvent(newEventMap map[string]string, attendeesEmails []string) (calendar.Event, error) {
	// Getting the authenticated calendar service
	srv, err := calendarAuth.GetCalendarService()

	if err != nil {
		log.Fatalf("Error: %v", err)
		panic(err)
	}
	// creating an array of attendees with the emails given as a parameter
	fmt.Println(len(attendeesEmails))
	attendees := make([]*calendar.EventAttendee, len(attendeesEmails), len(attendeesEmails))
	for _, v := range attendeesEmails {
		attendee := &calendar.EventAttendee{Email: v}
		attendees = append(attendees, attendee)
	}
	fmt.Println(len(attendees))
	// creating new event
	newEvent := &calendar.Event{
		Summary:     newEventMap["title"],
		Location:    newEventMap["location"],
		Description: newEventMap["description"],
		Start: &calendar.EventDateTime{
			DateTime: newEventMap["startDateTime"],
			TimeZone: "Egypt",
		},
		End: &calendar.EventDateTime{
			DateTime: newEventMap["endDateTime"],
			TimeZone: "Egypt",
		},
		Attendees: attendees,
		Organizer: &calendar.EventOrganizer{Email: newEventMap["organizerEmail"]},
	}

	calendarID := newEventMap["calenderID"]

	eventInsertCall := srv.Events.Insert(calendarID, newEvent)
	// send notification to attendees by email
	eventInsertCall.SendNotifications(true)
	// send request
	event, err := eventInsertCall.Do()

	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
		return *event, err
	}
	fmt.Println("event added  ", event)
	fmt.Println("event ID = ", event.Id)
	return *event, nil

}

//CreateAdvancedLabCalendar creates a new calendar for Advanced computer lab course
func CreateAdvancedLabCalendar() (calendar.Calendar, error) {
	// Getting the authenticated calendar service
	srv, err := calendarAuth.GetCalendarService()
	if err != nil {
		log.Fatalf("Error: %v", err)
		panic(err)
	}

	// creating new calendar
	newCalendar := &calendar.Calendar{
		// Id:          "b3cjs1oc7ql5jaecm98hv3fke0@group.calendar.google.com",
		Summary:     "Advanced Computer Lab Calendar",
		Description: "Calendar to add any events related to Advanced computer Lab course",
	}
	// Inserting new calendar
	calendar, err := srv.Calendars.Insert(newCalendar).Do()
	if err != nil {
		log.Fatalf("Unable to create calendar. %v\n", err)
		return *calendar, err
	}
	fmt.Println("calendar created and calendar id = ", calendar.Id)
	return *calendar, nil
}

//DeleteCalendar deletes a calendar with a specific ID
func DeleteCalendar(calendarID string) error {
	// Getting the authenticated calendar service
	srv, err := calendarAuth.GetCalendarService()
	if err != nil {
		log.Fatalf("Error: %v", err)
		panic(err)
	}
	// sending delete request to delete the calendar with ID = calendarID
	err = srv.Calendars.Delete(calendarID).Do()
	if err != nil {
		log.Fatalf("Unable to delete calendar. %v\n", err)
		return err
	}
	fmt.Println("calendar with ID: &i  deleted", calendarID)
	return nil

}

//UpdateEvent it updates a specific event
func UpdateEvent(calendarID, eventID string, newAttendees []string, deletedAttendees map[string]string, updatedEventMap map[string]string) (calendar.Event, error) {
	// Getting the authenticated calendar service
	srv, err := calendarAuth.GetCalendarService()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
	// Retrieve the event from the api
	event, err := srv.Events.Get(calendarID, eventID).Do()

	// updating the event fields with the new values that the user want to update
	for key, value := range updatedEventMap {
		key = strings.ToLower(key)
		switch key {
		case "title":
			event.Summary = value
		case "description":
			event.Description = value
		case "startDateTime":
			event.Start = &calendar.EventDateTime{
				DateTime: value,
				TimeZone: "Egypt",
			}
		case "endDateTime":
			event.End = &calendar.EventDateTime{
				DateTime: value,
				TimeZone: "Egypt",
			}
		case "location":
			event.Location = value
		case "organizerEmail":
			event.Organizer = &calendar.EventOrganizer{Email: value}
		default:
			return calendar.Event{}, fmt.Errorf("No field with this name: %s", key)
		}
	}

	if deletedAttendees != nil || newAttendees != nil {
		// get event attendees
		eventAttendees := event.Attendees

		// remove the attendees that user want to remove from event attendees
		if deletedAttendees != nil && len(deletedAttendees) > 0 {
			temp := []*calendar.EventAttendee{}
			for _, value := range eventAttendees {
				if _, ok := deletedAttendees[value.Email]; !ok {
					temp = append(temp, value)
				}
			}

			eventAttendees = temp

		}
		// add new Attendees to event Attendees
		if newAttendees != nil && len(newAttendees) > 0 {
			for _, v := range newAttendees {
				newEventAttendee := calendar.EventAttendee{
					Email: v,
				}
				eventAttendees = append(eventAttendees, &newEventAttendee)
			}
		}
		// update event attendees field with the new array
		event.Attendees = eventAttendees
	}

	event, err = srv.Events.Update(calendarID, eventID, event).Do()

	if err != nil {
		fmt.Println("Unable to update event. ", err)
		return calendar.Event{}, err
	}

	fmt.Println("Event Updated", event)
	return *event, nil
}

//ListEvents list all events in a specific calendar
func ListEvents(calendarID string) (string, []*calendar.Event, error) {
	// Getting the authenticated calendar service
	srv, err := calendarAuth.GetCalendarService()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	result, err := srv.Events.List(calendarID).Fields("items(description,id,summary)", "summary", "nextPageToken").Do()

	if err != nil {
		fmt.Println("Unable to retrieve calendar events list ", err)
		return "", make([]*calendar.Event, 0, 0), err
	}

	events := result.Items
	// fmt.Println(result.Summary)

	// ***pagination
	// for result.NextPageToken != "" {
	// 	result, err := srv.Events.List(calendarID).Fields("items(description,id,summary)", "summary", "nextPageToken").Do()
	// 	events = append(events[:], result.Items[:])
	// }

	// for _, event := range events {
	// 	fmt.Println("event id: " + event.Id + " and event summary: " + event.Summary)
	// }
	return result.Summary, events, nil
}

//GetEvent gets an event
func GetEvent(calendarID, eventID string) (calendar.Event, error) {
	// Getting the authenticated calendar service
	srv, err := calendarAuth.GetCalendarService()
	if err != nil {
		log.Fatalf("Error: %v", err)
		panic(err)
	}

	event, err := srv.Events.Get(calendarID, eventID).Do()
	if err != nil {
		fmt.Println("Unable to get event  ", err)
		return calendar.Event{}, fmt.Errorf("No Event with this ID: "+eventID+" ", err)
	}

	fmt.Println("event: ", event)
	fmt.Println("---------------------------------------")
	for _, v := range event.Attendees {
		fmt.Println(v)
	}
	return *event, nil
}

//DeleteEvent deletes an event
func DeleteEvent(calendarID, eventID string) error {
	// Getting the authenticated calendar service
	srv, err := calendarAuth.GetCalendarService()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	err = srv.Events.Delete(calendarID, eventID).Do()
	if err != nil {
		fmt.Println("Unable to delete event ", err)
		return err
	}
	fmt.Println("event with ID: " + eventID + " is deleted")
	return nil
}
