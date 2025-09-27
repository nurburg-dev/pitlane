package dbrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/entities"
)

type ActivityRunRepository interface {
	GetNextActivityRun(ctx context.Context) (*entities.DBActivityRun, error)
	GetActivityRunHistory(ctx context.Context, workflowRunId string) ([]entities.DBActivityRun, error)
	CreateActivityRun(ctx context.Context, activityRun *entities.DBActivityRun) error
	ChangeActivityRunStatus(ctx context.Context, activityRunID string, status entities.ActivityStatus) error
	GetActivityRun(ctx context.Context, activityRunID string) (*entities.DBActivityRun, error)
}

type PGActivityRunRepository struct {
	tx     pgx.Tx
	mapper *db.RowMapper
}

func NewPGActivityRunRepository(tx pgx.Tx) *PGActivityRunRepository {
	return &PGActivityRunRepository{
		tx:     tx,
		mapper: db.NewRowMapper(),
	}
}

func (r *PGActivityRunRepository) GetNextActivityRun(ctx context.Context) (*entities.DBActivityRun, error) {
	query := `
		SELECT id, activity_name, workflow_run_id, errorMessage, input, output,
			   status, retry_status, scheduled_at, created_at, updated_at
		FROM activity_runs
		WHERE status = @status
		ORDER BY scheduled_at DESC
		LIMIT 1
	`

	args := map[string]interface{}{
		"status": entities.ActivityStatusPending,
	}

	row := r.tx.QueryRow(ctx, query, pgx.NamedArgs(args))

	var activityRun entities.DBActivityRun
	err := r.mapper.ScanRow(row, &activityRun)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &activityRun, nil
}

func (r *PGActivityRunRepository) GetActivityRunHistory(
	ctx context.Context,
	workflowRunId string,
) ([]entities.DBActivityRun, error) {
	query := `
		SELECT id, activity_name, workflow_run_id, errorMessage, input, output,
			   status, retry_status, scheduled_at, created_at, updated_at
		FROM activity_runs
		WHERE workflow_run_id = @workflow_run_id
		ORDER BY created_at ASC
	`

	args := map[string]interface{}{
		"workflow_run_id": workflowRunId,
	}

	rows, err := r.tx.Query(ctx, query, pgx.NamedArgs(args))
	if err != nil {
		return nil, err
	}

	var activities []entities.DBActivityRun
	err = r.mapper.ScanRows(rows, &activities)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *PGActivityRunRepository) CreateActivityRun(ctx context.Context, activityRun *entities.DBActivityRun) error {
	query := `
		INSERT INTO activity_runs (id, activity_name, workflow_run_id, errorMessage, input, output,
								  status, retry_status, scheduled_at, created_at, updated_at)
		VALUES (@id, @activity_name, @workflow_run_id, @error_message, @input, @output,
				@status, @retry_status, @scheduled_at, @created_at, @updated_at)
	`

	args := map[string]interface{}{
		"id":              activityRun.ID,
		"activity_name":   activityRun.ActivityName,
		"workflow_run_id": activityRun.WorkflowRunID,
		"error_message":   activityRun.ErrorMessage,
		"input":           activityRun.Input,
		"output":          activityRun.Output,
		"status":          activityRun.Status,
		"retry_status":    activityRun.RetryStatus,
		"scheduled_at":    activityRun.ScheduledAt,
		"created_at":      activityRun.CreatedAt,
		"updated_at":      activityRun.UpdatedAt,
	}

	_, err := r.tx.Exec(ctx, query, pgx.NamedArgs(args))
	return err
}

func (r *PGActivityRunRepository) ChangeActivityRunStatus(
	ctx context.Context,
	activityRunID string,
	status entities.ActivityStatus,
) error {
	query := `
		UPDATE activity_runs
		SET status = @status, updated_at = NOW()
		WHERE id = @id
	`

	args := map[string]interface{}{
		"id":     activityRunID,
		"status": status,
	}

	_, err := r.tx.Exec(ctx, query, pgx.NamedArgs(args))
	return err
}

func (r *PGActivityRunRepository) GetActivityRun(
	ctx context.Context,
	activityRunID string,
) (*entities.DBActivityRun, error) {
	query := `
		SELECT id, activity_name, workflow_run_id, errorMessage, input, output,
			   status, retry_status, scheduled_at, created_at, updated_at
		FROM activity_runs
		WHERE id = @id
	`

	args := map[string]interface{}{
		"id": activityRunID,
	}

	row := r.tx.QueryRow(ctx, query, pgx.NamedArgs(args))

	var activityRun entities.DBActivityRun
	err := r.mapper.ScanRow(row, &activityRun)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &activityRun, nil
}
