package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Here we define the roles a user can have.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// Claims is the payload of JWTs.
type Claims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

// NewClaims creates a new Claims object for the identified user. Additional fields can be
// set after creating this Claim.
func NewClaims(subject string, roles []string, now time.Time, expires time.Duration) Claims {

	c := Claims{
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(expires).Unix(),
		},
	}

	return c
}
