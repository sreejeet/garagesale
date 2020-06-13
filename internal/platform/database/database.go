package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The drier being used
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
