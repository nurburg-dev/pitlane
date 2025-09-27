package pitlane

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/dbrepo"
	"github.com/nurburg-dev/pitlane/internal/entities"
)

type WorkflowWorker struct {
	pgPool *pgxpool.Pool
}

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
				ww.markWorkflowRunErrored(ctx, wfRun.ID, err.Error())
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
		// TODO: store activityHistory in a ordered list.
		// TODO: implement a ExecuteActivity function which looks for activity reponse in history ordered list
		// TODO: if there is a history mismatch in history order list then return error
		// TODO: if execution has reached history ordered list end then return and specific error to schedule activity
	})
}
