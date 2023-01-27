// Package product contains product related CRUD functionality.
package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"log"
	"time"
)

var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

const (
	insertOp = "insert"
	deleteOp = "delete"
)

const (
	iq = `INSERT INTO products
				(id, name, price, date_created, date_updated) 
				VALUES($1, $2, $3, $4, $5)`
	lq = `INSERT INTO products_logs(product_id, operation) VALUES ($1, $2)`
	dq = `DELETE FROM products WHERE id = $1`
)

func Create(ctx context.Context, db *sqlx.DB, np NewProduct) (Product, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "business.data.product.create")
	defer span.End()

	p := Product{
		ID:          uuid.NewString(),
		Name:        np.Name,
		Price:       np.Price,
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	tx, err := db.Begin()
	if err != nil {
		return Product{}, fmt.Errorf("beginning transaction: %w", err)
	}
	_, err = tx.ExecContext(ctx, iq, p.ID, p.Name, p.Price, p.DateCreated, p.DateUpdated)
	_, err = tx.ExecContext(ctx, lq, p.ID, insertOp)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return Product{}, fmt.Errorf("rollback transaction: %w", err)
		}
		return Product{}, fmt.Errorf("inserting log product: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return Product{}, fmt.Errorf("committing transaction: %w", err)
	}
	return p, nil
}

// Delete removes the product identified by a given ID.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {
	ctx, span := otel.Tracer("service").Start(ctx, "business.data.product.delete")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	_, err = db.ExecContext(ctx, dq, id)
	_, err = tx.ExecContext(ctx, lq, id, deleteOp)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("rollback transaction: %w", err)
		}
		return fmt.Errorf("inserting log product: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// Search finds products.
func Search(ctx context.Context, db *elasticsearch.Client, query string) ([]SearchProduct, error) {
	ctx, span := otel.Tracer("service").Start(ctx, "business.data.product.search")
	defer span.End()
	var buf bytes.Buffer
	searchQuery := map[string]any{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": query,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := db.Search(
		db.Search.WithContext(ctx),
		db.Search.WithIndex("products"),
		db.Search.WithBody(&buf),
		db.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("searching product: %w", err)
	}
	defer res.Body.Close()
	var products map[string]any
	var searchProducts []SearchProduct
	if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("decoding products: %w", err)
	}
	for _, hit := range products["hits"].(map[string]any)["hits"].([]any) {
		var searchProduct SearchProduct
		prod := hit.(map[string]interface{})["_source"]
		prodByte, _ := json.Marshal(prod)
		_ = json.Unmarshal(prodByte, &searchProduct)
		searchProducts = append(searchProducts, searchProduct)
	}
	return searchProducts, nil
}
