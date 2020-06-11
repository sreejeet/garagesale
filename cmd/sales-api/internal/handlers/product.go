package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sreejeet/garagesale/internal/product"
)

// Products define all handlers for products. It also holds
// the application state needed by the handler methods.
type Products struct {
	DB *sqlx.DB
}

// List is an http handler for returning
// a json list of products.
func (p *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		log.Printf("error listing products: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshalling (converting) product slice to json string
	data, err := json.Marshal(list)
	if err != nil {
		log.Print("Error parsing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		log.Print("Error writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
