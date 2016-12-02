package routesHandlers

import (
	"encoding/json"
	"models"
	"net/http"

	calendar "google.golang.org/api/calendar/v3"
)

type errorObj struct {
	Resource string `json:"resource"`
	Message  string `json:"message"`
}

type errorsJSONObj struct {
	Status  int64      `json:"Status"`
	Message string     `json:"message"`
	Errors  []errorObj `json:"errors"`
}

type successListJSONObj struct {
	Status        int64             `json:"Status"`
	Message       string            `json:"message"`
	Page          int64             `json:"page"`
	CalendarTitle string            `json:"calendar"`
	Results       []*calendar.Event `json:"results"`
}
type successSingleJSONObj struct {
	Status  int64               `json:"Status"`
	Message string              `json:"message"`
	Results []map[string]string `json:"results"`
}

type successEventJSONObj struct {
	Status  int64          `json:"Status"`
	Message string         `json:"message"`
	Results calendar.Event `json:"results"`
}

type successJSONObj struct {
	Status  int64             `json:"Status"`
	Message string            `json:"message"`
	Results map[string]string `json:"results"`
}

type successTaskJSONObj struct {
	Status  int64       `json:"Status"`
	Message string      `json:"message"`
	Results models.Task `json:"results"`
}

type successUserJSONObj struct {
	Status  int64       `json:"Status"`
	Message string      `json:"message"`
	Results models.User `json:"results"`
}

type successTasksListJSONObj struct {
	Status  int64             `json:"Status"`
	Message string            `json:"message"`
	Results []models.TaskList `json:"results"`
}

//WriteJSON Writes the JSON equivilant for data into ResponseWriter w
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
