package database

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The drier being used
	"go.opencensus.io/trace"
)

// Config struct holds the required database parameters
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

// Open workds as an abstraction to open a database conn
func Open(cfg Config) (*sqlx.DB, error) {

	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	query := make(url.Values)
	query.Set("sslmode", sslMode)
	query.Set("timezone", "utc")

	url := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: query.Encode(),
	}

	return sqlx.Open("postgres", url.String())
}

// StatusCheck returns an error in case the databse is not working as expected
// or nil if everythign is fine.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	ctx, span := trace.StartSpan(ctx, "platform.DB.StatusCheck")
	defer span.End()

	// Avoid running a ping request for health checks as it may return a false positive.
	// Run a qurey instead to make sure the database is accepting and honoring requests.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
