package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"pull-request-review/config"
	"pull-request-review/internal/infrastructure/adapters/logger"
)

type Database struct {
	pool   *pgxpool.Pool
	cfg    config.DatabaseConfig
	logger logger.Logger
}

func NewDatabase(cfg config.DatabaseConfig, log logger.Logger) *Database {
	return &Database{
		cfg:    cfg,
		logger: log,
	}
}

func (d *Database) Connect(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(d.cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}

	poolConfig.MaxConns = d.cfg.MaxConns
	poolConfig.MinConns = d.cfg.MinConns
	poolConfig.MaxConnLifetime = d.cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = d.cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = d.cfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	d.pool = pool
	d.logger.Info("Database connection established")

	return nil
}

func (d *Database) Ping(ctx context.Context) error {
	if d.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return d.pool.Ping(ctx)
}

func (d *Database) Close() {
	if d.pool != nil {
		d.logger.Info("Closing database connection pool")
		d.pool.Close()
	}
}

func (d *Database) GetPool() *pgxpool.Pool {
	return d.pool
}