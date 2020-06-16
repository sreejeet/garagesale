package product

import (
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
func Retrieve(db *sqlx.DB) (*Product, error) {
	var p Product

	const query = `SELECT * FROM products WHERE product_id = $1`
	if err := db.Get(&p, query); err != nil {
		return nil, errors.Wrap(err, "selecting one product")
	}
	return &p, nil
}
