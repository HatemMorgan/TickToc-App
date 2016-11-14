package routes

type errorObj struct {
	Resource string `json:"resource"`
	Message  string `json:"message"`
}

type errorsJSONObj struct {
	Status  string     `json:"Status"`
	Message string     `json:"message"`
	Errors  []errorObj `json:"errors"`
}

type successJSONObj struct {
	Status  string        `json:"Status"`
	Message string        `json:"message"`
	Page    int64         `json:"page"`
	Results []interface{} `json:"results"`
}
