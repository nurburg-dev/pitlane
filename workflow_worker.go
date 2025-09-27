package pitlane

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/dbrepo"
	"github.com/nurburg-dev/pitlane/internal/entities"
	"github.com/nurburg-dev/pitlane/internal/utils"
)

type WorkflowWorker struct {
	pgPool *pgxpool.Pool
}
type contextKey string

const (
	activityHistoryKey      contextKey = "activityHistory"
	activityRepoKey         contextKey = "activityRepo"
	workflowRunIDKey        contextKey = "workflowRunID"
	activityHistoryIndexKey contextKey = "activityHistoryIndex"
)

var (
	ErrActivityHistoryMismatch  = errors.New("activity history mismatch")
	ErrActivitySchedulingNeeded = errors.New("activity scheduling needed")
)

func (ww *WorkflowWorker) Start() {
	go func() {
		for {
			time.Sleep(time.Second)
			ctx := context.TODO()
			wfRun, err := ww.pickNextWorkflow(ctx)
			if err != nil {
				fmt.Printf("WARN: error while consuming workflow - %s", err.Error())
			}
			if wfRun == nil {
				continue
			}
			fmt.Printf("INFO: found run %s", wfRun.ID)
			err = ww.executeWorkflow(ctx, wfRun)
			if err != nil {
				if !errors.Is(err, ErrActivitySchedulingNeeded) {
					if markErr := ww.markWorkflowRunErrored(ctx, wfRun.ID, err.Error()); markErr != nil {
						fmt.Printf("WARN: failed to mark workflow as errored: %s", markErr.Error())
					}
				}
			}
		}
	}()
}

func (ww *WorkflowWorker) pickNextWorkflow(ctx context.Context) (*entities.DBWorkflowRun, error) {
	var pickedWf *entities.DBWorkflowRun = nil
	err := db.ExecuteTx(ctx, ww.pgPool, func(ctx context.Context, tx pgx.Tx) error {
		wfRepo := dbrepo.NewPGWorkflowRepository(tx)
		wf, err := wfRepo.GetNextWorkflowRun(ctx)
		pickedWf = wf
		if err != nil {
			return err
		}
		if wf == nil {
			return nil
		}
		return wfRepo.ChangeWorkflowRunStatus(ctx, wf.ID, entities.WorkflowStatusExecuting)
	})
	if err != nil {
		return nil, err
	}
	return pickedWf, nil
}

func (ww *WorkflowWorker) markWorkflowRunErrored(ctx context.Context, wfRunID, errorMessage string) error {
	return db.ExecuteTx(ctx, ww.pgPool, func(ctx context.Context, tx pgx.Tx) error {
		wfRepo := dbrepo.NewPGWorkflowRepository(tx)
		return wfRepo.MarkWorkflowRunErrored(ctx, wfRunID, errorMessage)
	})
}

func (ww *WorkflowWorker) executeWorkflow(ctx context.Context, wfRun *entities.DBWorkflowRun) error {
	return db.ExecuteTx(ctx, ww.pgPool, func(ctx context.Context, tx pgx.Tx) error {
		activityRepo := dbrepo.NewPGActivityRunRepository(tx)
		activityHistory, err := activityRepo.GetActivityRunHistory(ctx, wfRun.ID)
		if err != nil {
			return err
		}

		orderedActivityHistory := make([]*entities.DBActivityRun, len(activityHistory))
		for i := range activityHistory {
			orderedActivityHistory[i] = &activityHistory[i]
		}

		workflowStore := GetWorkflowStore()
		workflowFunc, exists := workflowStore[wfRun.WorkflowName]
		if !exists {
			return fmt.Errorf("workflow function %s not found in registry", wfRun.WorkflowName)
		}

		ctx = context.WithValue(ctx, activityHistoryKey, orderedActivityHistory)
		ctx = context.WithValue(ctx, activityRepoKey, activityRepo)
		ctx = context.WithValue(ctx, workflowRunIDKey, wfRun.ID)
		ctx = context.WithValue(ctx, activityHistoryIndexKey, &[]int{0})

		var args []interface{}
		if unmarshalErr := json.Unmarshal(wfRun.Input, &args); unmarshalErr != nil {
			return fmt.Errorf("failed to unmarshal workflow input: %w", unmarshalErr)
		}

		// Prepend context to args for reflection call
		reflectArgs := make([]interface{}, len(args)+1)
		reflectArgs[0] = ctx
		copy(reflectArgs[1:], args)

		results, err := utils.InvokeFunction(workflowFunc, reflectArgs...)
		if err != nil {
			return fmt.Errorf("failed to invoke workflow function: %w", err)
		}

		// Check if the function returned an error (second return value)
		if len(results) >= 2 && !results[1].IsNil() {
			if err, ok := results[1].Interface().(error); ok {
				if errors.Is(err, ErrActivitySchedulingNeeded) {
					wfRepo := dbrepo.NewPGWorkflowRepository(tx)
					if statusErr := wfRepo.ChangeWorkflowRunStatus(
						ctx,
						wfRun.ID,
						entities.WorkflowStatusWaitingActivity,
					); statusErr != nil {
						return fmt.Errorf("failed to change workflow status to waiting for activity: %w", statusErr)
					}
				}
				return err
			}
		}

		// If no error, mark workflow as finished
		wfRepo := dbrepo.NewPGWorkflowRepository(tx)
		return wfRepo.ChangeWorkflowRunStatus(ctx, wfRun.ID, entities.WorkflowStatusFinished)
	})
}
