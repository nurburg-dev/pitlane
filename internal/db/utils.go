package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TableExists(ctx context.Context, pool *pgxpool.Pool, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`
	err := pool.QueryRow(ctx, query, tableName).Scan(&exists)
	return exists, err
}
