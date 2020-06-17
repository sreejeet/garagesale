package product

import "time"

// Product is an individial item that can be sold.
type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewProduct type is expected from clients when creating a product.
type NewProduct struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}
