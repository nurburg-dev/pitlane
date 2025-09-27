package db

import (
	"context"
	"math/rand"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

const idSuffixLength = 8

func GenerateReadableID() string {
	prefix := petname.Generate(2, "-")
	randomSuffix := generateRandomString(idSuffixLength)
	return prefix + "-" + randomSuffix
}
