package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	//set up configuration
	config.MaxConns = 10
	config.MaxConns = 2
	config.MaxConnLifetime = 10 * time.Minute

	return pgxpool.NewWithConfig(context.Background(), config)
}
