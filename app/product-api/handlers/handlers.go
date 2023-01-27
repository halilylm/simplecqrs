// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/halilylm/simplecqrs/business/mid"
	"github.com/halilylm/simplecqrs/foundation/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, es *elasticsearch.Client, db *sqlx.DB) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	c := check{
		build: build,
		db:    db,
		es:    es,
	}
	app.Handle(http.MethodGet, "/v1/health", c.health)

	p := productHandlers{
		db: db,
		es: es,
	}

	app.Handle(http.MethodGet, "/v1/products", p.search)
	app.Handle(http.MethodPost, "/v1/products", p.create)
	app.Handle(http.MethodDelete, "/v1/products/:id", p.delete)

	return app
}
