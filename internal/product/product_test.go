package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sreejeet/garagesale/internal/platform/auth"
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

	// Create a claims object with some random UUID for testing
	claims := auth.NewClaims(
		"718ffbea-f4a1-4667-8ae3-b349da52675e",
		[]string{auth.RoleAdmin, auth.RoleUser},
		now, time.Hour,
	)

	p0, err := product.Create(ctx, db, claims, newP, now)
	if err != nil {
		t.Fatalf("creating product: %s", err)
	}

	p1, err := product.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting product: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched product not same as created product:\n%s", diff)
	}

	update := product.UpdateProduct{
		Name: tests.StringPointer("Updated Name"),
		Cost: tests.IntPointer(51),
	}
	updatedTime := time.Date(2020, time.January, 1, 1, 1, 1, 0, time.UTC)

	if err := product.Update(ctx, db, claims, p0.ID, update, updatedTime); err != nil {
		t.Fatalf("updating product p0: %s", err)
	}

	saved, err := product.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	// Check specified fields were updated. Make a copy of the original product
	// and change just the fields we expect then diff it with what was saved.
	want := *p0
	want.Name = "Updated Name"
	want.Cost = 51
	want.DateUpdated = updatedTime

	if diff := cmp.Diff(want, *saved); diff != "" {
		t.Fatalf("updated record did not match:\n%s", diff)
	}

	// Check if product delete works
	if err := product.Delete(ctx, db, p0.ID); err != nil {
		t.Fatalf("deleting product: %v", err)
	}

	_, err = product.Retrieve(ctx, db, p0.ID)
	if err == nil {
		t.Fatalf("should not be able to retrieve deleted product")
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
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
