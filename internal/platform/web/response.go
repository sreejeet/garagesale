package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond function encodes the data into json
// and writes it into the response writer.
func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {

	res, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	w.Header().Set("Content_Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(res); err != nil {
		return err
	}

	return nil
}

// RespondError is used to send error responses to the client.
func RespondError(w http.ResponseWriter, err error) error {

	// Check if type is of *Error, that means it was an expected error
	// and it may contain a specific error code that must be used instead of 500.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(w, er, webErr.Status); err != nil {
			return err
		}
	} else {
		// The error was unexpected, send Internal Server Error code.
		er := ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		}
		if err := Respond(w, er, http.StatusInternalServerError); err != nil {
			return err
		}
	}

	return nil
}
