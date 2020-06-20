package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Handler is the func signature used by all handlers in this service
type Handler func(http.ResponseWriter, *http.Request) error

// App will be the entry point to our REST API.
// It will control the context of each request.
type App struct {
	log *log.Logger
	mux *chi.Mux
}

// NewApp is a contructor for REST API App
func NewApp(log *log.Logger) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
	}
}

// Handle associates a handlerfunc with an HTTP method and URL pattern.
// This converts our custom handler to the standard lib Handler type.
// It captures errors and returns them to the client in a consistent manner.
func (a *App) Handle(method, url string, h Handler) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {

			// Logging to our logs
			a.log.Printf("Error: %v+", err)

			// Respond with the error to the client
			if err := RespondError(w, err) err !=nil{
				a.log.Printf("Error: %v+", err)
			}
			Respond(w, res, http.StatusInternalServerError)
		}
	}

	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
