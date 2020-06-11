package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The drier being used
)

// Open workds as an abstraction to open a database conn
func Open() (*sqlx.DB, error) {
	query := make(url.Values)
	query.Set("sslmode", "disable")
	query.Set("timezone", "utc")

	url := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}
}
