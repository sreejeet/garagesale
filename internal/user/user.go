package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Create is used to create a new user.
func Create(ctx context.Context, db *sqlx.DB, n NewUser, now time.Time) (*User, error) {

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
