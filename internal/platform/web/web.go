package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

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
func (a *App) Handle(method, url string, h http.HandlerFunc) {
	a.mux.MethodFunc(method, url, h)
}

// ServeHTTP implements the http.Handler interface
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
