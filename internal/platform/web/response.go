package web

import (
	"encoding/json"
	"net/http"
)

// Respond function encodes the data into json
// and writes it into the response writer.
func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {

	res, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	w.Header().Set("Content_Type", "application/json; charset=utf-8")
	if _, err := w.Write(res); err != nil {
		return err
	}

	return nil
}
