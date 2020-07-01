package web

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	var u struct {
		Name string `validate:"required"`
	}

	// Test the decoding function with missing arguments
	body := strings.NewReader(`{}`)

	r := httptest.NewRequest("POST", "/", body)
	// The expected outcome is an error.
	// The test must fail if we do not get an error from the decoder.
	err := Decode(r, &u)
	if err == nil {
		t.Errorf("Decode with missing arguments should return an error but returned nil")
	}

	t.Log(err)
}
