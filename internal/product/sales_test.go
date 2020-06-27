package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/sreejeet/garagesale/internal/product"
	"github.com/sreejeet/garagesale/internal/tests"
)

func TestSales(t *testing.T) {

	db, teardown := tests.NewUnit(t)
	defer teardown()

	now := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	ctx := context.Background()

	// Create two products as seed data.
	newPuzzles := product.NewProduct{
		Name:     "Puzzles",
		Cost:     25,
		Quantity: 6,
	}
	puzzles, err := product.Create(ctx, db, newPuzzles, now)
	if err != nil {
		t.Fatalf("creating product: %s", err)
	}

	newToys := product.NewProduct{
		Name:     "Toys",
		Cost:     40,
		Quantity: 3,
	}
	toys, err := product.Create(ctx, db, newToys, now)
	if err != nil {
		t.Fatalf("creating product: %s", err)
	}

	{ // Add and list sales

		ns := product.NewSale{
			Quantity: 3,
			Paid:     70,
		}

		s, err := product.AddSale(ctx, db, ns, puzzles.ID, now)
		if err != nil {
			t.Fatalf("creating sale: %s", err)
		}

		// Puzzles should show the one sale added above.
		sales, err := product.ListSales(ctx, db, puzzles.ID)
		if err != nil {
			t.Fatalf("listing sales: %s", err)
		}
		if exp, got := 1, len(sales); exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}
		if exp, got := s.ID, sales[0].ID; exp != got {
			t.Fatalf("expected sale ID %v, got %v", exp, got)
		}

		// Toys should have 0 sales.
		sales, err = product.ListSales(ctx, db, toys.ID)
		if err != nil {
			t.Fatalf("listing sales: %s", err)
		}
		if exp, got := 0, len(sales); exp != got {
			t.Fatalf("expected sale list size %v, got %v", exp, got)
		}
	}
}
