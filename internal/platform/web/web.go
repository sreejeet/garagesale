package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values are stored or retrieved.
const KeyValues ctxKey = 1

// Values struct stores key information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

// Handler is the func signature used by all handlers in this service
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// App will be the entry point to our REST API.
// It will control the context of each request.
type App struct {
	log *log.Logger
	mux *chi.Mux
	mw  []Middleware
}

// NewApp is a contructor for REST API App
func NewApp(log *log.Logger, mw ...Middleware) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
		mw:  mw,
	}
}

// Handle associates a handlerfunc with an HTTP method and URL pattern.
// This converts our custom handler to the standard lib Handler type.
// It captures errors and returns them to the client in a consistent manner.
func (a *App) Handle(method, url string, h Handler, mw ...Middleware) {

	// Wrap with specific middleware provided
	h = wrapMiddleware(mw, h)

	// Wrap with other application middlewares
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {

		// Create a Values struct to record state for the request. Store the
		// address in the request's context so it is sent down the call chain.
		v := Values{
			Start: time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		// Run and catch any exeption from the handler chain.
		if err := h(ctx, w, r); err != nil {
			// Logging to our logs
			a.log.Printf("Unexpected err: %v+", err)
		}
	}

	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
