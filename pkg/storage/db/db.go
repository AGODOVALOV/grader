// Package db provides database functionality.
package db

import (
	"context"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/storage/db/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RepoDB is a database connection pool.
type RepoDB struct {
	Pool *pgxpool.Pool
}

// NewRepoDB creates a new database connection pool.
func NewRepoDB(ctx context.Context, cfg *config.Config) (*RepoDB, error) {
	db, err := postgres.NewPostgresDB(ctx, &cfg.DB)
	if err != nil {
		return nil, err
	}
	return &RepoDB{Pool: db}, nil
}
