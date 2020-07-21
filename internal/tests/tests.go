package tests

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sreejeet/garagesale/internal/platform/auth"
	"github.com/sreejeet/garagesale/internal/platform/database"
	"github.com/sreejeet/garagesale/internal/platform/database/databasetest"
	"github.com/sreejeet/garagesale/internal/schema"
	"github.com/sreejeet/garagesale/internal/user"
)

// NewUnit creates a test database inside a container and creates the reqired table structue.
// In case of a failiour, it will call Fatal on the testing.T parameter.
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
	maxAttempts := 300
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
		t.Fatalf("Failed starting database after %d seconds: %s", maxAttempts/10, pingError)
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

// Test owns state for running and shutting down tests.
type Test struct {
	DB            *sqlx.DB
	Log           *log.Logger
	Authenticator *auth.Authenticator

	t       *testing.T
	cleanup func()
}

// New creates a database, seeds it, constructs an authenticator.
func New(t *testing.T) *Test {
	t.Helper()

	// Initialize and seed database. Store the cleanup function call later.
	db, cleanup := NewUnit(t)

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	// Create the logger to use.
	logger := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// Create RSA keys to enable authentication in our service.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Build an authenticator using this static key.
	kid := "4754d86b-7a6d-4df5-9c65-224741361492"
	kf := auth.NewSimpleKeyLookupFunc(kid, key.Public().(*rsa.PublicKey))
	authenticator, err := auth.NewAuthenticator(key, kid, "RS256", kf)
	if err != nil {
		t.Fatal(err)
	}

	return &Test{
		DB:            db,
		Log:           logger,
		Authenticator: authenticator,
		t:             t,
		cleanup:       cleanup,
	}
}

// Teardown releases any resources used for the test.
func (test *Test) Teardown() {
	test.cleanup()
}

// Token generates an authenticated token for a user.
func (test *Test) Token(email, pass string) string {
	test.t.Helper()

	claims, err := user.Authenticate(
		context.Background(), test.DB, time.Now(),
		email, pass,
	)
	if err != nil {
		test.t.Fatal(err)
	}

	tkn, err := test.Authenticator.GenerateToken(claims)
	if err != nil {
		test.t.Fatal(err)
	}

	return tkn
}

// StringPointer is a helper function to return a pointer to a string.
// We do not need this outside testing so it is declared here.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper function to return a pointer to an int.
// We do not need this outside testing so it is declared here.
func IntPointer(i int) *int {
	return &i
}
