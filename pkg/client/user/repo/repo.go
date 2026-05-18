// Package repo provides a user repository.
package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/storage/db"
)

// Repo is a user repository.
type Repo struct {
	Queries *Queries
	DB      *db.RepoDB
}

// NewRepo creates a new user repository.
func NewRepo(ctx context.Context, cfg *config.Config) (*Repo, error) {
	repoDB, err := db.NewRepoDB(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Queries: New(repoDB.Pool),
		DB:      repoDB,
	}, nil
}

// ExecTx executes a transaction.
func (r *Repo) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := r.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.New(fmt.Sprintf("tx err: %v, rb err: %v", err, rbErr))
		}
		return err
	}
	return tx.Commit(ctx)
}
