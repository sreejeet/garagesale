package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sreejeet/garagesale/internal/platform/web"
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
func (p *Products) List(w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(p.db)
	if err != nil {
		return errors.Wrap(err, "Error listing products")
	}

	// Using the web.Respond helper to return json
	return web.Respond(w, list, http.StatusOK)
}

// Retrieve is used to get a single product based on its ID from the URL parameter.
func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.db, id)
	if err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrap(err, "Error finding product")
		}
	}

	// Using the web.Respond helper to return json
	return web.Respond(w, prod, http.StatusOK)
}

// Create is used to create a new product from the body of a request.
// The created product is sent back to the client
// in conformance to the RESTful architecture.
func (p *Products) Create(w http.ResponseWriter, r *http.Request) error {

	var newProd product.NewProduct

	// Decoding response body to NewProduct struct
	if err := web.Decode(r, &newProd); err != nil {
		return errors.Wrap(err, "Error decoding product")
	}

	// Creating product in database
	prod, err := product.Create(p.db, newProd, time.Now())
	if err != nil {
		return errors.Wrap(err, "Error creating product")
	}

	// Using the web.Respond helper to return json
	return web.Respond(w, &prod, http.StatusOK)
}
