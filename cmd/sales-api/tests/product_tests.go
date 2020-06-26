package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sreejeet/garagesale/cmd/sales-api/internal/handlers"
	"github.com/sreejeet/garagesale/internal/schema"
	"github.com/sreejeet/garagesale/internal/tests"
)

// ProductTests is used for passing dependencies for tests and also simplify
// adding subtests.
type ProductTests struct {
	app http.Handler
}

func TestProducts(t *testing.T) {

	// Create a new unit for testing
	db, teardown := tests.NewUnit(t)
	defer teardown()

	// Seed the database int he above created unit
	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	// Set where and how to log
	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tests := ProductTests{app: handlers.API(db, log)}

	t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

// List tests the listing of products from the API
func (p *ProductTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/products", nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected http status code %v, got %v", http.StatusOK, resp.Code)
	}

	// Using a slice of empty interface to decode the response.
	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("Decoding list of products: %s", err)
	}

	// The exprected list of products.
	// This list must be exactly the same as the products defined in the database seeding function.
	want := []map[string]interface{}{
		{
			"id":           "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":         "Comic Books",
			"cost":         float64(50),
			"quantity":     float64(42),
			"date_created": "2019-01-01T00:00:01.000001Z",
			"date_updated": "2019-01-01T00:00:01.000001Z",
		},
		{
			"id":           "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":         "McDonalds Toys",
			"cost":         float64(75),
			"quantity":     float64(120),
			"date_created": "2019-01-01T00:00:02.000001Z",
			"date_updated": "2019-01-01T00:00:02.000001Z",
		},
	}

	// Check if the API response and expected results are the same or not.
	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

// ProductCRUD test will be used to perform all CRUD operations of the API
func (p *ProductTests) ProductCRUD(t *testing.T) {

	var created map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"product0","cost":55,"quantity":6}`)

		req := httptest.NewRequest("POST", "/v1/products", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Fatalf("Posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("Decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("Expected non-empty product id")
		}
		if created["date_created"] == "" || created["date_created"] == nil {
			t.Fatal("Expected non-empty product date_created")
		}
		if created["date_updated"] == "" || created["date_updated"] == nil {
			t.Fatal("Expected non-empty product date_updated")
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
			"name":         "product0",
			"cost":         float64(55),
			"quantity":     float64(6),
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("Retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("Decoding: %s", err)
		}

		// Fetched product should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("Retrieved product should match created. Diff:\n%s", diff)
		}
	}
}