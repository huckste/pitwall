package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	connString := os.Getenv("DATABASE_URL")
	return pgxpool.New(ctx, connString)
}
