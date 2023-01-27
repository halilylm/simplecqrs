// Package database provides support for access the database
package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"go.opentelemetry.io/otel"
	"net/url"
)

// Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

// Open knows how to open a database connection based on the config.
func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	fmt.Println(u.String())
	return sqlx.Open("postgres", u.String())
}

// StatusCheck returns nil if it can successfully talk to the database.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	ctx, span := otel.Tracer("service").Start(ctx, "foundation.database.StatusCheck")
	defer span.End()

	const q = `SELECT 1`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
