package txrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/nurburg-dev/pitlane"
	"github.com/nurburg-dev/pitlane/internal/db"
)

type WorkflowRepository interface {
	GetNextWorkflowRun(ctx context.Context) (*pitlane.DBWorkflowRun, error)
	GetWorkflow(ctx context.Context, name string) (*pitlane.DBWorkflow, error)
	UpsertWorkflow(ctx context.Context, workflow *pitlane.DBWorkflow) error
	CreateWorkflowRun(ctx context.Context, workflowRun *pitlane.DBWorkflowRun) error
	ChangeWorkflowRunStatus(ctx context.Context, workflowRunID string, status pitlane.WorkflowStatus) error
}

type PGWorkflowRepository struct {
	tx     pgx.Tx
	mapper *db.RowMapper
}

func NewPGWorkflowRepository(tx pgx.Tx) *PGWorkflowRepository {
	return &PGWorkflowRepository{
		tx:     tx,
		mapper: db.NewRowMapper(),
	}
}

func (r *PGWorkflowRepository) GetNextWorkflowRun(ctx context.Context) (*pitlane.DBWorkflowRun, error) {
	query := `
		SELECT id, input, workflow_name, status, scheduled_at, created_at, updated_at
		FROM workflow_runs
		WHERE status = @status
		ORDER BY scheduled_at DESC
		LIMIT 1
	`

	args := map[string]interface{}{
		"status": pitlane.WorkflowStatusPending,
	}

	row := r.tx.QueryRow(ctx, query, pgx.NamedArgs(args))

	var workflowRun pitlane.DBWorkflowRun
	err := r.mapper.ScanRow(row, &workflowRun)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &workflowRun, nil
}

func (r *PGWorkflowRepository) GetWorkflow(ctx context.Context, name string) (*pitlane.DBWorkflow, error) {
	query := `
		SELECT name, created_at, updated_at
		FROM workflows
		WHERE name = @name
	`

	args := map[string]interface{}{
		"name": name,
	}

	row := r.tx.QueryRow(ctx, query, pgx.NamedArgs(args))

	var workflow pitlane.DBWorkflow
	err := r.mapper.ScanRow(row, &workflow)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &workflow, nil
}

func (r *PGWorkflowRepository) UpsertWorkflow(ctx context.Context, workflow *pitlane.DBWorkflow) error {
	query := `
		INSERT INTO workflows (name, created_at, updated_at)
		VALUES (@name, @created_at, @updated_at)
		ON CONFLICT (name) DO UPDATE SET
			updated_at = @updated_at
	`

	args := map[string]interface{}{
		"name":       workflow.Name,
		"created_at": workflow.CreatedAt,
		"updated_at": workflow.UpdatedAt,
	}

	_, err := r.tx.Exec(ctx, query, pgx.NamedArgs(args))
	return err
}

func (r *PGWorkflowRepository) ChangeWorkflowRunStatus(ctx context.Context, workflowRunID string, status pitlane.WorkflowStatus) error {
	query := `
		UPDATE workflow_runs
		SET status = @status, updated_at = NOW()
		WHERE id = @id
	`

	args := map[string]interface{}{
		"id":     workflowRunID,
		"status": status,
	}

	_, err := r.tx.Exec(ctx, query, pgx.NamedArgs(args))
	return err
}

func (r *PGWorkflowRepository) CreateWorkflowRun(ctx context.Context, workflowRun *pitlane.DBWorkflowRun) error {
	query := `
		INSERT INTO workflow_runs (id, input, workflow_name, status, scheduled_at, created_at, updated_at)
		VALUES (@id, @input, @workflow_name, @status, @scheduled_at, @created_at, @updated_at)
	`

	args := map[string]interface{}{
		"id":            workflowRun.ID,
		"input":         workflowRun.Input,
		"workflow_name": workflowRun.WorkflowName,
		"status":        workflowRun.Status,
		"scheduled_at":  workflowRun.ScheduledAt,
		"created_at":    workflowRun.CreatedAt,
		"updated_at":    workflowRun.UpdatedAt,
	}

	_, err := r.tx.Exec(ctx, query, pgx.NamedArgs(args))
	return err
}
