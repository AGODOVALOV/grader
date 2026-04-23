package db

import (
	"context"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/storage/db/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoDB struct {
	Pool *pgxpool.Pool
}

func NewRepoDB(ctx context.Context, cfg *config.Config) (*RepoDB, error) {
	db, err := postgres.NewPostgresDB(ctx, &cfg.DB)
	if err != nil {
		return nil, err
	}
	return &RepoDB{Pool: db}, nil
}
