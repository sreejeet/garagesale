package web

import (
	"encoding/json"
	"net/http"
)

// Decode function decodes a json response into the provided value.
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}
	return nil
}
