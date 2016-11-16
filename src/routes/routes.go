package routes

import (
	"log"
	"net/http"
	"net/http/httptest"
	"routesHandlers"

	cors "github.com/heppu/simple-cors"
)

//Routing handles all Routings
func Routing(addr string) error {
	// HandleFuncs

	mux := http.NewServeMux()

	mux.HandleFunc("/welcome", withLog(routesHandlers.HandleWelcome))
	mux.HandleFunc("/chat/event", withLog(routesHandlers.HandleChat))
	mux.HandleFunc("/events/list", withLog(routesHandlers.EventListHandler))
	mux.HandleFunc("/events", withLog(routesHandlers.EventHandler))
	mux.HandleFunc("/tasks", withLog(routesHandlers.TaskHandler))
	mux.HandleFunc("/users", withLog(routesHandlers.UsersHandler))
	mux.HandleFunc("/", withLog(routesHandlers.Handle))
	// Start the server
	return http.ListenAndServe(addr, cors.CORS(mux))
}

//WithLog Wraps HandlerFuncs to log requests to Stdout
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
