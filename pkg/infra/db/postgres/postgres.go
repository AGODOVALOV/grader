package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/AGODOVALOV/grader/pkg/infra/db/postgres/config"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

func NewPostgresDB(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&timezone=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode, cfg.TimeZone)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Pool.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Pool.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.Pool.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		logger.Z(ctx).Error(ctx, "create Postgres connection", err.Error(), map[string]string{
			"host":    cfg.Host,
			"port":    fmt.Sprintf("%d", cfg.Port),
			"db":      cfg.DBName,
			"user":    cfg.User,
			"sslmode": cfg.SSLMode,
		})
		return nil, err
	}

	logger.Z(ctx).Info(ctx, "create Postgres connection", "ping ok", map[string]string{
		"host":    cfg.Host,
		"port":    fmt.Sprintf("%d", cfg.Port),
		"db":      cfg.DBName,
		"user":    cfg.User,
		"sslmode": cfg.SSLMode,
	})

	return db, nil
}
