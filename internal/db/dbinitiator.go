package db

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBInitator interface {
	Init(ctx context.Context) error
}

//go:embed pgschema.sql
var schemaSQL string

type PGInitiator struct {
	pool *pgxpool.Pool
}

func NewPGInitiator(pool *pgxpool.Pool) *PGInitiator {
	return &PGInitiator{
		pool: pool,
	}
}

func (p *PGInitiator) Init(ctx context.Context) error {
	_, err := p.pool.Exec(ctx, schemaSQL)
	return err
}
