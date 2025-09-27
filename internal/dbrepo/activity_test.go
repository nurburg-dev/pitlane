package dbrepo_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/dbrepo"
	"github.com/nurburg-dev/pitlane/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGActivityRunRepository_CreateAndGetNextActivityRun(t *testing.T) {
	ctx := context.Background()

	// Get connection from pool
	conn, err := testContainer.GetPool().Acquire(ctx)
	require.NoError(t, err)
	defer conn.Release()

	// Start transaction
	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Create workflow repository and insert test data
	workflowRepo := dbrepo.NewPGWorkflowRepository(tx)

	// Create test workflow
	now := time.Now()
	workflow := &entities.DBWorkflow{
		Name:      "test-workflow",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = workflowRepo.UpsertWorkflow(ctx, workflow)
	require.NoError(t, err)

	// Create test workflow run
	workflowRunID := db.GenerateReadableID()
	workflowRun := &entities.DBWorkflowRun{
		ID:           workflowRunID,
		Input:        json.RawMessage(`{}`),
		WorkflowName: "test-workflow",
		Status:       entities.WorkflowStatusPending,
		ScheduledAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = workflowRepo.CreateWorkflowRun(ctx, workflowRun)
	require.NoError(t, err)

	// Create activity repository
	repo := dbrepo.NewPGActivityRunRepository(tx)

	// Test data
	activityRun := &entities.DBActivityRun{
		ID:            db.GenerateReadableID(),
		ActivityName:  "test-activity",
		WorkflowRunID: workflowRunID,
		Input:         json.RawMessage(`{"test": "data"}`),
		Status:        entities.ActivityStatusPending,
		ScheduledAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Test CreateActivityRun
	err = repo.CreateActivityRun(ctx, activityRun)
	require.NoError(t, err)

	// Test GetNextActivityRun
	nextActivity, err := repo.GetNextActivityRun(ctx)
	require.NoError(t, err)
	require.NotNil(t, nextActivity)

	assert.Equal(t, activityRun.ID, nextActivity.ID)
	assert.Equal(t, activityRun.ActivityName, nextActivity.ActivityName)
	assert.Equal(t, activityRun.WorkflowRunID, nextActivity.WorkflowRunID)
	assert.Equal(t, activityRun.Status, nextActivity.Status)
	assert.JSONEq(t, string(activityRun.Input), string(nextActivity.Input))

	// Test ChangeActivityRunStatus
	err = repo.ChangeActivityRunStatus(ctx, activityRun.ID, entities.ActivityStatusExecuting)
	require.NoError(t, err)

	// Verify status changed by getting the next activity run (should be nil since it's executing)
	nextActivityAfterUpdate, err := repo.GetNextActivityRun(ctx)
	require.NoError(t, err)
	require.Nil(t, nextActivityAfterUpdate)
}
