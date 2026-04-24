// Package repo provides user repository.
package repo

import (
	"context"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/storage/db"
)

// Repo is a user repository.
type Repo struct {
	Queries *Queries
	db      *db.RepoDB
}

// NewRepo creates a new user repository.
func NewRepo(ctx context.Context, cfg *config.Config) (*Repo, error) {
	repoDB, err := db.NewRepoDB(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Queries: New(repoDB.Pool),
		db:      repoDB,
	}, nil
}
