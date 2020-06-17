package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
	data, err := json.Marshal(list)
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
	data, err := json.Marshal(prod)
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
