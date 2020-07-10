package handlers

import (
	"context"
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
func (p *Products) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(ctx, p.db)
	if err != nil {
		return errors.Wrap(err, "Error listing products")
	}

	// Using the web.Respond helper to return json
	return web.Respond(ctx, w, list, http.StatusOK)
}

// Retrieve is used to get a single product based on its ID from the URL parameter.
func (p *Products) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(ctx, p.db, id)
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
	return web.Respond(ctx, w, prod, http.StatusOK)
}

// Create is used to create a new product from the body of a request.
// The created product is sent back to the client
// in conformance to the RESTful architecture.
func (p *Products) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var newProd product.NewProduct

	// Decoding response body to NewProduct struct
	if err := web.Decode(r, &newProd); err != nil {
		return errors.Wrap(err, "Error decoding product")
	}

	// Creating product in database
	prod, err := product.Create(ctx, p.db, newProd, time.Now())
	if err != nil {
		return errors.Wrap(err, "Error creating product")
	}

	// Using the web.Respond helper to return json
	return web.Respond(ctx, w, &prod, http.StatusCreated)
}

// AddSale records a new sale transaction for a specific product.
// It takes a NewSale object in json from and returns the added record to the caller.
func (p *Products) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(ctx, p.db, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(ctx, w, sale, http.StatusCreated)
}

// ListSales lists all sales for a specific product.
func (p *Products) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(ctx, p.db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Update takes the product id from the url and updates the fields that have been provided to it.
func (p *Products) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(ctx, p.db, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a specific product from the database based on the give id.
func (p *Products) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	if err := product.Delete(ctx, p.db, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
