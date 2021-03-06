package product

import "time"

// Product is an individial item that can be sold.
type Product struct {
	ID          string    `db:"product_id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Sold        int       `db:"sold" json:"sold"`
	Revenue     int       `db:"revenue" json:"revenue"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewProduct type is expected from clients when creating a product.
type NewProduct struct {
	Name     string `json:"name" validate:"required"`
	Cost     int    `json:"cost" validate:"gte=0"`
	Quantity int    `json:"quantity" validate:"gte=1"`
}

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateProduct struct {
	Name     *string `json:"name"`
	Cost     *int    `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int    `json:"quantity" validate:"omitempty,gte=1"`
}

// Sale type denotes a single sale transaction of a product.
// Quantity is the number of items of a product were sold in this transaction.
// Paid is the cumulative amount that was paid for this transaction
type Sale struct {
	ID          string    `db:"sale_id" json:"id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	Quantity    int       `db:"quantity" json:"quantity" validate:"gte=0"`
	Paid        int       `db:"paid" json:"paid" validate:"gte=0"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
}

// NewSale is the form for recording a transaction.
type NewSale struct {
	Quantity int `json:"quantity"`
	Paid     int `json:"paid"`
}
