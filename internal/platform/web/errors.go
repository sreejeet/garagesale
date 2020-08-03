package web

import "github.com/pkg/errors"

// ErrorResponse is used as the default response type for any API errors.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error type passes the error with a specific web status code.
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

// Error implements the error interface. It uses the default error
// of the wrapped error. This is what will be visible on the services' logs.
func (err *Error) Error() string {
	return err.Err.Error()
}

// FieldError indicates an erro with a specific field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// NewRequestError simply wraps an error with a status code.
// Only to be used by service handlers in case of known errors.
func NewRequestError(err error, status int) error {
	return &Error{err, status, nil}
}

// shutdown is a type used to help with the graceful termination of the service.
// shutdown type is used for the graceful termination of the service
// in case a critical or unexpected error has occoured.
type shutdown struct {
	Message string
}

// Error is the implementation of the error interface on a shutdown type instance
func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// IsShutdown checks to see if the shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
