package handlers

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/halilylm/simplecqrs/business/data/product"
	"github.com/halilylm/simplecqrs/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"net/http"
)

type productHandlers struct {
	db *sqlx.DB
	es *elasticsearch.Client
}

func (h *productHandlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := otel.Tracer("service").Start(ctx, "app.product-api.handlers.productHandlers.create")
	defer span.End()

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return fmt.Errorf("decoding new product: %w", err)
	}
	prod, err := product.Create(ctx, h.db, np)
	if err != nil {
		return fmt.Errorf("creating new product: %w", err)
	}

	return web.Respond(ctx, w, prod, http.StatusCreated)
}

func (h *productHandlers) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := otel.Tracer("service").Start(ctx, "app.product-api.handlers.productHandlers.delete")
	defer span.End()

	params := web.Params(r)
	if err := product.Delete(ctx, h.db, params["id"]); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("product[%s]: %w", params["id"], err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *productHandlers) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := otel.Tracer("service").Start(ctx, "app.product-api.handlers.productHandlers.search")
	defer span.End()

	products, err := product.Search(ctx, h.es, r.URL.Query().Get("name"))
	if err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("query[%s]: %w", r.URL.Query().Get("name"), err)
		}
	}

	return web.Respond(ctx, w, products, http.StatusOK)
}
