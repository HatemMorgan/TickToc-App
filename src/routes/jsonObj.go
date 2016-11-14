package routes

import (
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
	Status  int64               `json:"Status"`
	Message string              `json:"message"`
	Page    int64               `json:"page"`
	Results []map[string]string `json:"results"`
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
