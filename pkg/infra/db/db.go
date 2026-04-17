package db

import (
	"context"
	"database/sql"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/infra/db/postgres"
)

type RepoDB struct {
	db *sql.DB
}

func NewRepoDB(ctx context.Context, cfg *config.Config) (*RepoDB, error) {
	db, err := postgres.NewPostgresDB(ctx, &cfg.DB)

	if err != nil {
		return nil, err
	}

	return &RepoDB{db: db}, nil
}
