package handlers

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/halilylm/simplecqrs/foundation/database"
	"github.com/halilylm/simplecqrs/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"net/http"
)

type check struct {
	build string
	db    *sqlx.DB
	es    *elasticsearch.Client
}

// If the services are not ready we will tell the client and use a 502
// status.
func (c *check) health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := otel.Tracer("service").Start(ctx, "app.product-api.handlers.health")
	defer span.End()

	status := "ok"
	statusCode := http.StatusOK
	if err := database.StatusCheck(ctx, c.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusBadGateway
	}

	health := struct {
		Version string `json:"version"`
		Status  string `json:"status"`
	}{
		Version: c.build,
		Status:  status,
	}

	return web.Respond(ctx, w, health, statusCode)
}
