// Package postgres provides PostgreSQL database connection.
package postgres

import (
	"context"
	"fmt"
	"net/url"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/db/postgres/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // #nosec G104
)

// NewPostgresDB creates a new PostgreSQL database connection pool.
func NewPostgresDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.DBName,
	}

	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	q.Set("TimeZone", cfg.TimeZone)
	u.RawQuery = q.Encode()

	dsn := u.String()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = int32(cfg.Pool.MaxOpenConns) // #nosec G115
	poolConfig.MinConns = int32(cfg.Pool.MaxIdleConns) // #nosec G115
	poolConfig.MaxConnLifetime = cfg.Pool.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.Pool.ConnMaxIdleTime
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Z(ctx).Error(ctx, "create Postgres connection", err.Error(), map[string]string{
			"host":    cfg.Host,
			"port":    fmt.Sprintf("%d", cfg.Port),
			"db":      cfg.DBName,
			"user":    cfg.User,
			"sslmode": cfg.SSLMode,
		})
		return nil, err
	}
	return pool, nil
}
