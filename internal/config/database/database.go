package database

import (
	"context"
	"crud/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(dbConfig config.DatabaseConfig) (*pgxpool.Pool, error) {
	connectionString := dbConfig.ToConnectionString()
	return pgxpool.New(context.Background(), connectionString)
}
