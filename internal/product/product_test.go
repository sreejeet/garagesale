package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sreejeet/garagesale/internal/product"
	"github.com/sreejeet/garagesale/internal/schema"
	"github.com/sreejeet/garagesale/internal/tests"
)

func TestProducts(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	newP := product.NewProduct{
		Name:     "Bite my shiny metal as - Bender B Rodríguez",
		Cost:     2999,
		Quantity: 1,
	}
	now := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	p0, err := product.Create(ctx, db, newP, now)
	if err != nil {
		t.Fatalf("Creating product: %s", err)
	}

	p1, err := product.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("Getting product: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("Fetched product not same as created product:\n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := product.List(context.Background(), db)
	if err != nil {
		t.Fatalf("Listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("Expected product list size %v, got %v", exp, got)
	}
}
