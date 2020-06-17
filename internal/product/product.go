package product

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// List retrieves all products from the database
func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}
	const query = `SELECT * FROM products`
	if err := db.Select(&products, query); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}
	return products, nil
}

// Retrieve is used to get a single product based on its ID from the URL parameter.
func Retrieve(db *sqlx.DB, id string) (*Product, error) {
	var prod Product

	const query = `SELECT * FROM products WHERE product_id = $1`
	if err := db.Get(&prod, query, id); err != nil {
		return nil, errors.Wrap(err, "selecting one product")
	}
	return &prod, nil
}

// Create creates a new product int the database and return the created product.
func Create(db *sqlx.DB, newProd NewProduct, now time.Time) (*Product, error) {
	prod := Product{
		ID:          uuid.New().String(),
		Name:        newProd.Name,
		Cost:        newProd.Cost,
		Quantity:    newProd.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const query = `INSERT INTO products
		(product_id, name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(query,
		prod.ID, prod.Name,
		prod.Cost, prod.Quantity,
		prod.DateCreated, prod.DateUpdated)

	if err != nil {
		return nil, errors.Wrap(err, "Creating new product")
	}

	return &prod, nil
}
