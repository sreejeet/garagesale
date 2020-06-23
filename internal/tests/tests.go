package tests

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sreejeet/garagesale/internal/platform/database"
	"github.com/sreejeet/garagesale/internal/platform/database/databasetest"
	"github.com/sreejeet/garagesale/internal/schema"
)

// NewUnit creates a test database inside a container and creates the reqired table structue.
// In case of a failiour, it will call Fatal on the testing.T.
func NewUnit(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := databasetest.StartContainer(t)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})

	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	// We will wait for the databse to be ready.
	// We will ping the database every 100ms till we dont get an error.
	t.Log("Waiting for database to be ready")
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Check if we have a successful connection
	if pingError != nil {
		databasetest.DumpContainerLogs(t, c)
		databasetest.StopContainer(t, c)
		t.Fatalf("Failed starting database: %s", pingError)
	}

	// Perform schema migration
	if err := schema.Migrate(db); err != nil {
		databasetest.StopContainer(t, c)
		t.Fatalf("Migration failed %s", err)
	}

	// teardown function is called after the caller is done with the test.
	teardown := func() {
		t.Helper()
		db.Close()
		databasetest.StopContainer(t, c)
	}

	return db, teardown
}
