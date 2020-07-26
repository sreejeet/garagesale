package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// AddSale records a single sale transaction for a product.
func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale, productID string, now time.Time) (*Sale, error) {

	ctx, span := trace.StartSpan(ctx, "internal.product.AddSale")
	defer span.End()

	s := Sale{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Quantity:    ns.Quantity,
		Paid:        ns.Paid,
		DateCreated: now,
	}

	const q = `INSERT INTO sales
		(sale_id, product_id, quantity, paid, date_created)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := db.ExecContext(ctx, q,
		s.ID, s.ProductID, s.Quantity,
		s.Paid, s.DateCreated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating sale")
	}

	return &s, nil
}

// ListSales lists all sale transactions for a product.
func ListSales(ctx context.Context, db *sqlx.DB, productID string) ([]Sale, error) {

	ctx, span := trace.StartSpan(ctx, "internal.product.ListSales")
	defer span.End()

	sales := []Sale{}

	const q = `SELECT * FROM sales WHERE product_id = $1`
	if err := db.SelectContext(ctx, &sales, q, productID); err != nil {
		return nil, errors.Wrap(err, "listing sales")
	}

	return sales, nil
}
