// Package elastic provides support for access the database
package elastic

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

type Config struct {
	URL string
}

func Open(cfg Config) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.URL},
	},
	)
	if err != nil {
		return nil, errors.Wrap(err, "elastic.NewClient")
	}

	return client, nil
}
