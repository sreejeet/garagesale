package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sreejeet/garagesale/internal/platform/web"
)

// API constructs an app instance with all application routes defined
func API(db *sqlx.DB, log *log.Logger) http.Handler {

	app := web.NewApp(log)

	p := Products{
		db:  db,
		log: log,
	}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	return app
}
