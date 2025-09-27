package pitlane

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/dbrepo"
	"github.com/nurburg-dev/pitlane/internal/entities"
	"github.com/nurburg-dev/pitlane/internal/utils"
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
		err := we.initializeDB(ctx)
		if err != nil {
			return nil, err
		}
	}
	return we, nil
}

func (we *WorkflowEngine) initializeDB(ctx context.Context) error {
	dbInitiator := db.NewPGInitiator(we.pgPool)
	return dbInitiator.Init(ctx)
}

func (we *WorkflowEngine) InvokeWorkflow(ctx context.Context, workflowFunction any, args ...any) (string, error) {
	workflowFuncName, err := utils.GetFunctionName(workflowFunction)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow function name: %w", err)
	}

	if _, exist := GetWorkflowStore()[workflowFuncName]; !exist {
		return "", fmt.Errorf("workflow %s not registered", workflowFuncName)
	}

	if err2 := utils.ValidateArgs(workflowFunction, args...); err2 != nil {
		return "", err2
	}

	inputBytes, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workflow input: %w", err)
	}
	now := time.Now()

	tx, err := we.pgPool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	workflowRepo := dbrepo.NewPGWorkflowRepository(tx)

	err = workflowRepo.UpsertWorkflow(
		ctx,
		&entities.DBWorkflow{
			Name:      workflowFuncName,
			CreatedAt: now,
			UpdatedAt: now,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upsert workflow: %w", err)
	}

	workflowRunID := db.GenerateReadableID()
	workflowRun := &entities.DBWorkflowRun{
		ID:           workflowRunID,
		Input:        inputBytes,
		WorkflowName: workflowFuncName,
		Status:       entities.WorkflowStatusPending,
		ScheduledAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = workflowRepo.CreateWorkflowRun(ctx, workflowRun)
	if err != nil {
		return "", fmt.Errorf("failed to create workflow run: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return workflowRunID, nil
}
