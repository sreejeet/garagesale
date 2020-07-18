package auth

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// KeyLookupFunc is used to map a JWT key id (kid) to the corresponding public key.
// It is a requirement for creating an Authenticator.
//
// * Private keys should be rotated. During the transition period, tokens
// signed with the old and new keys can coexist by looking up the correct
// public key by key id (kid).
//
// * Key-id-to-public-key resolution is usually accomplished via a public JWKS
// endpoint. See https://auth0.com/docs/jwks for more details.
type KeyLookupFunc func(kid string) (*rsa.PublicKey, error)

// NewSimpleKeyLookupFunc is used to make a key lookup function that returns a single public key based on the active key id.
func NewSimpleKeyLookupFunc(activeKID string, publicKey *rsa.PublicKey) KeyLookupFunc {

	// Return KeyLookupFunc
	f := func(kid string) (*rsa.PublicKey, error) {
		if activeKID != kid {
			return nil, fmt.Errorf("unrecognized key id %q", kid)
		}
		return publicKey, nil
	}

	return f
}

// Authenticator authenticates clients. It can generate and
// also recreate calaims by parsing tokens
type Authenticator struct {
	privateKey       *rsa.PrivateKey
	activeKID        string
	algorithm        string
	pubKeyLookupFunc KeyLookupFunc
	parser           *jwt.Parser
}

// NewAuthenticator is a factory function for creating Authenticator instances.
func NewAuthenticator(privateKey *rsa.PrivateKey, activeKID, algorithm string, publicKeyLookupFunc KeyLookupFunc) (*Authenticator, error) {

	// Validate privided parameters for potential errors.
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}
	if activeKID == "" {
		return nil, errors.New("active key id cannot be blank")
	}
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}
	if publicKeyLookupFunc == nil {
		return nil, errors.New("public key function cannot be nil")
	}

	// Create a parser that can be used to validate the
	// signing algorithm used to avoid a critical bug in JWT.
	// Refer: https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := Authenticator{
		privateKey:       privateKey,
		activeKID:        activeKID,
		algorithm:        algorithm,
		pubKeyLookupFunc: publicKeyLookupFunc,
		parser:           &parser,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Authenticator) GenerateToken(claims Claims) (string, error) {

	// Use the signing algorithm used by the Authenticator
	method := jwt.GetSigningMethod(a.algorithm)

	// Create token using singing method and claims.
	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = a.activeKID

	// Sign the token using the Authenticator's private key.
	// Return error if failed.
	str, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return str, nil
}
