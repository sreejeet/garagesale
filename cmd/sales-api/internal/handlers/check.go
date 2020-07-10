package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sreejeet/garagesale/internal/platform/database"
	"github.com/sreejeet/garagesale/internal/platform/web"
)

// Check holds health check points for the service
type Check struct {
	db *sqlx.DB
}

// Health provides a high level overview of the service status
// and lets the client know if and why it is [not]capable of taking requests
func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var health struct {
		DBStatus string `json:"db_status"`
	}

	// Check if the database is ready.
	if err := database.StatusCheck(ctx, c.db); err != nil {

		// If the database is not ready we will respond to the client with a 500 status.
		// Returning the error as it is will result in an unhandled exception,
		// send a json web response insetead.
		health.DBStatus = "The database is not accpeting requests at the moment."
		return web.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.DBStatus = "ok"
	return web.Respond(ctx, w, health, http.StatusOK)
}
