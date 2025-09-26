package pitlane

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurburg-dev/pitlane/internal/db"
)

type WorkflowEngine struct {
	pgPool *pgxpool.Pool
}

func NewWorkflowEngine(ctx context.Context, config *EngineConfig) (*WorkflowEngine, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBConfig.Username,
		config.DBConfig.Password,
		config.DBConfig.Host,
		config.DBConfig.Port,
		config.DBConfig.Database,
	)

	pgPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	we := &WorkflowEngine{
		pgPool: pgPool,
	}
	if config.InitDB {
		err := we.init(ctx)
		if err != nil {
			return nil, err
		}
	}
	return we, nil
}

func (we *WorkflowEngine) init(ctx context.Context) error {
	dbInitiator := db.NewPGInitiator(we.pgPool)
	return dbInitiator.Init(ctx)
}
