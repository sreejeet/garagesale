package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sreejeet/garagesale/internal/mid"
	"github.com/sreejeet/garagesale/internal/platform/auth"
	"github.com/sreejeet/garagesale/internal/platform/web"
)

// API constructs an app instance with all application routes defined
func API(db *sqlx.DB, log *log.Logger, authenticator *auth.Authenticator) http.Handler {

	// App holds all the routes as well as the middleware chain
	app := web.NewApp(
		log,
		mid.Logger(log),
		mid.Errors(log),
		mid.Metrics(),
	)

	{
		c := Check{db: db}

		// Health check route
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		// User authentication routes
		u := Users{db: db, authenticator: authenticator}
		app.Handle(http.MethodGet, "/v1/users/token", u.Token)
	}

	{
		p := Products{
			db:  db,
			log: log,
		}

		// Product specific routes
		app.Handle(http.MethodGet, "/v1/products", p.List)
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
		app.Handle(http.MethodPost, "/v1/products", p.Create)
		app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
		app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

		// Sale specific routes
		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)
	}

	return app
}
