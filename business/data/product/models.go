package product

import "time"

// Product is an item we sell
type Product struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Price       int       `db:"price" json:"price"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewProduct is what we require from clients when adding a Product.
type NewProduct struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required,gte=0"`
}

// SearchProduct is what we get from Elasticsearch
type SearchProduct struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Operation string `json:"operation"`
	Price     int    `json:"price"`
}
