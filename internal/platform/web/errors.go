package web

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
