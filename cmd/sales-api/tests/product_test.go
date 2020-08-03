package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sreejeet/garagesale/cmd/sales-api/internal/handlers"
	"github.com/sreejeet/garagesale/internal/tests"
)

// ProductTests is used for passing dependencies for tests and also simplify
// adding subtests.
type ProductTests struct {
	app        http.Handler
	adminToken string
}

func TestProducts(t *testing.T) {

	// Create a new unit for testing
	test := tests.New(t)
	defer test.Teardown()

	shutdown := make(chan os.Signal, 1)
	tests := ProductTests{
		app: handlers.API(
			shutdown,
			test.DB,
			test.Log,
			test.Authenticator,
		),
		adminToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

// List tests the listing of products from the API
func (p *ProductTests) List(t *testing.T) {

	req := httptest.NewRequest("GET", "/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+p.adminToken)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected http status code %v, got %v", http.StatusOK, resp.Code)
	}

	// Using a slice of empty interface to decode the response.
	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding list of products: %s", err)
	}

	// The exprected list of products.
	// This list must be exactly the same as the products defined in the database seeding function.
	want := []map[string]interface{}{
		{
			"id":           "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":         "Comic Books",
			"cost":         float64(50),
			"quantity":     float64(42),
			"revenue":      float64(350),
			"sold":         float64(7),
			"user_id":      "00000000-0000-0000-0000-000000000000",
			"date_created": "2019-01-01T00:00:01.000001Z",
			"date_updated": "2019-01-01T00:00:01.000001Z",
		},
		{
			"id":           "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":         "McDonalds Toys",
			"cost":         float64(75),
			"quantity":     float64(120),
			"revenue":      float64(225),
			"sold":         float64(3),
			"user_id":      "00000000-0000-0000-0000-000000000000",
			"date_created": "2019-01-01T00:00:02.000001Z",
			"date_updated": "2019-01-01T00:00:02.000001Z",
		},
	}

	// Check if the API response and expected results are the same or not.
	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("response did not match expected. Diff:\n%s", diff)
	}
}

// ProductCRUD test will be used to perform all CRUD operations of the API
func (p *ProductTests) ProductCRUD(t *testing.T) {

	var created map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"product0","cost":55,"quantity":6}`)

		req := httptest.NewRequest("POST", "/v1/products", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("expected non-empty product id")
		}
		if created["date_created"] == "" || created["date_created"] == nil {
			t.Fatal("expected non-empty product date_created")
		}
		if created["date_updated"] == "" || created["date_updated"] == nil {
			t.Fatal("expected non-empty product date_updated")
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
			"name":         "product0",
			"cost":         float64(55),
			"quantity":     float64(6),
			"sold":         float64(0),
			"revenue":      float64(0),
			"user_id":      tests.AdminID,
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// Fetched product should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("retrievedd product should match created. Diff:\n%s", diff)
		}
	}

	{ // UPDATE
		body := strings.NewReader(`{"name":"Updated Name","cost":20,"quantity":10}`)
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("PUT", url, body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusNoContent {
			t.Fatalf("updating: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve updated record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp = httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var updated map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": updated["date_updated"],
			"name":         "Updated Name",
			"cost":         float64(20),
			"quantity":     float64(10),
			"sold":         float64(0),
			"revenue":      float64(0),
			"user_id":      tests.AdminID,
		}

		// Updated product should match the one we created.
		if diff := cmp.Diff(want, updated); diff != "" {
			t.Fatalf("retrieved product should match created. Diff:\n%s", diff)
		}
	}

	{ // DELETE
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusNoContent {
			t.Fatalf("deleting: expected status code %v, got %v", http.StatusNoContent, resp.Code)
		}

		// Retrieve deleted record to be sure it worked.
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.adminToken)
		resp = httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusNotFound {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusNotFound, resp.Code)
		}
	}

}

// CreateRequiresFields tests the request decoder for proper validation checks
func (p *ProductTests) CreateRequiresFields(t *testing.T) {

	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/v1/products", body)

	req.Header.Set("Authorization", "Bearer "+p.adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %v, got %v", http.StatusBadRequest, resp.Code)
	}
}
