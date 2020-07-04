package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Custom errors for expected failing conditions
var (
	// Invalid UUID
	ErrInvalidID = errors.New("invalid ID")
	// Unable to find product based on UUID
	ErrNotFound = errors.New("product not found")
)

// List retrieves all products from the database
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {

	products := []Product{}
	const query = `SELECT
						p.*,
						COALESCE(SUM(s.quantity), 0) AS sold,
						COALESCE(SUM(s.paid), 0) AS revenue
					FROM products AS p
					LEFT JOIN sales AS s ON p.product_id = s.product_id
					GROUP BY p.product_id`

	if err := db.SelectContext(ctx, &products, query); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Retrieve is used to get a single product based on its ID from the URL parameter.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {

	// Check for invalid UUID
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var prod Product

	const query = `SELECT
						p.*,
						COALESCE(SUM(s.quantity), 0) AS sold,
						COALESCE(SUM(s.paid), 0) AS revenue
					FROM products AS p
					LEFT JOIN sales AS s ON p.product_id = s.product_id
					WHERE p.product_id = $1
					GROUP BY p.product_id`

	if err := db.GetContext(ctx, &prod, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "selecting one product")
	}

	return &prod, nil
}

// Create creates a new product int the database and return the created product.
func Create(ctx context.Context, db *sqlx.DB, newProd NewProduct, now time.Time) (*Product, error) {

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

	_, err := db.ExecContext(ctx, query,
		prod.ID, prod.Name,
		prod.Cost, prod.Quantity,
		prod.DateCreated, prod.DateUpdated)

	if err != nil {
		return nil, errors.Wrap(err, "Creating new product")
	}

	return &prod, nil
}

// Update modifies an existing product.
func Update(ctx context.Context, db *sqlx.DB, id string, update UpdateProduct, now time.Time) error {

	// Use the retrieve function to get the product to be updated.
	p, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	// Only update fields that have been passed as all fields are optional
	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Cost != nil {
		p.Cost = *update.Cost
	}
	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}
	p.DateUpdated = now

	const q = `UPDATE products SET
               "name" = $2,
               "cost" = $3,
               "quantity" = $4,
               "date_updated" = $5
               WHERE product_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		p.Name, p.Cost,
		p.Quantity, p.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}

// Delete removes products from the database based on the id provided.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM products WHERE product_id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting product %s", id)
	}

	return nil
}
