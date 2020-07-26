package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/crypto/bcrypt"

	"github.com/sreejeet/garagesale/internal/platform/auth"
)

var (
	// ErrAuthenticationFailure is used to indicate any kind of authentication error.
	ErrAuthenticationFailure = errors.New("Authentication failed")
)

// Create is used to create a new user.
func Create(ctx context.Context, db *sqlx.DB, n NewUser, now time.Time) (*User, error) {

	ctx, span := trace.StartSpan(ctx, "internal.user.Create")
	defer span.End()

	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		Roles:        n.Roles,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	// Insert the new user into the database.
	const q = `INSERT INTO users
               (user_id, name, email, password_hash, roles, date_created, date_updated)
               VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.ExecContext(
		ctx, q,
		u.ID, u.Name, u.Email,
		u.PasswordHash, u.Roles,
		u.DateCreated, u.DateUpdated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating new user")
	}

	return &u, nil
}

// Authenticate searches the database for a user based on the email, verifies the password
// and if this succeeds, a Claims object is returned. This Claims object is used to
// create a token for future authenctications.
func Authenticate(ctx context.Context, db *sqlx.DB, now time.Time, email, password string) (auth.Claims, error) {

	ctx, span := trace.StartSpan(ctx, "internal.user.Authenticate")
	defer span.End()

	// Query for email search
	const q = `SELECT * FROM users WHERE email = $1`

	var u User
	if err := db.GetContext(ctx, &u, q, email); err != nil {

		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}

		return auth.Claims{}, errors.Wrap(err, "selecting single user")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// Finally after the above checks have been successful, genereate user token.
	claims := auth.NewClaims(u.ID, u.Roles, now, time.Hour)
	return claims, nil
}
