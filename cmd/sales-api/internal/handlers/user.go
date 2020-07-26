package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sreejeet/garagesale/internal/platform/auth"
	"github.com/sreejeet/garagesale/internal/platform/web"
	"github.com/sreejeet/garagesale/internal/user"
	"go.opencensus.io/trace"
)

// Users holds handlers for dealing with user.
type Users struct {
	db            *sqlx.DB
	authenticator *auth.Authenticator
}

// Token creates an auth token for the user after authenticating themselves with an email and password.
func (u *Users) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Users.Token")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("web value missing from context")
	}

	// Get the Basic Auth credentials
	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	// Use the email and password to authenticate from the database.
	// If the user is authenticated, get the claims of the user.
	claims, err := user.Authenticate(ctx, u.db, v.Start, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}

	// Create a new token usingthe claims of the user returned from the databse.
	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = u.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
