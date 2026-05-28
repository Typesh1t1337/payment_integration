package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context, DBURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(DBURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 10
	cfg.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}