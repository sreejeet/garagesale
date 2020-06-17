package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/sreejeet/garagesale/internal/product"
)

// Products define all handlers for products. It also holds
// the application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// List is an http handler for returning
// a json list of products.
func (p *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.db)
	if err != nil {
		p.log.Println("Error listing products:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshalling (converting) product slice to json string
	data, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		p.log.Println("Error parsing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		p.log.Println("Error writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Retrieve is used to get a single product based on its ID from the URL parameter.
func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.db, id)
	if err != nil {
		p.log.Println("error finding product:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshalling (converting) product slice to json string
	data, err := json.MarshalIndent(prod, "", "    ")
	if err != nil {
		p.log.Println("Error parsing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		p.log.Println("Error writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Create is used to create a new product from the body of a request.
// The created product is sent back to the client
// in conformance to the RESTful architecture.
func (p *Products) Create(w http.ResponseWriter, r *http.Request) {

	var newProd product.NewProduct

	// Decoding response body to NewProduct struct
	if err := json.NewDecoder(r.Body).Decode(&newProd); err != nil {
		p.log.Println("Error decoding response:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prod, err := product.Create(p.db, newProd, time.Now())
	if err != nil {
		p.log.Println("Error creating product:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(prod, "", "    ")
	if err != nil {
		p.log.Println("Error marshalling result:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		p.log.Println("Error writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
