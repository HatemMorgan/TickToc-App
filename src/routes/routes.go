package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	// Third party packages

	cors "github.com/heppu/simple-cors"
)

//Routing handles all Routings
func Routing(addr string) error {
	// HandleFuncs
	mux := http.NewServeMux()

	mux.HandleFunc("/welcome", withLog(handleWelcome))
	mux.HandleFunc("/chat/event", withLog(handleChat))
	mux.HandleFunc("/", withLog(handle))
	// mux.HandleFunc("/events/list")
	mux.HandleFunc("/events", withLog(eventHandler))

	// Start the server
	return http.ListenAndServe(addr, cors.CORS(mux))
}

// withLog Wraps HandlerFuncs to log requests to Stdout
func withLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := httptest.NewRecorder()
		fn(c, r)
		log.Printf("[%d] %-4s %s\n", c.Code, r.Method, r.URL.Path)

		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		c.Body.WriteTo(w)
	}
}

// writeJSON Writes the JSON equivilant for data into ResponseWriter w
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
