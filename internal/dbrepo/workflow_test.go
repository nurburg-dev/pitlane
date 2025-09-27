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

func TestPGWorkflowRepository_UpsertAndGetWorkflow(t *testing.T) {
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

	// Create repository
	repo := dbrepo.NewPGWorkflowRepository(tx)

	// Test data
	now := time.Now()
	workflow := &entities.DBWorkflow{
		Name:      "test-workflow",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test UpsertWorkflow
	err = repo.UpsertWorkflow(ctx, workflow)
	require.NoError(t, err)

	// Test GetWorkflow
	retrievedWorkflow, err := repo.GetWorkflow(ctx, "test-workflow")
	require.NoError(t, err)
	require.NotNil(t, retrievedWorkflow)

	assert.Equal(t, workflow.Name, retrievedWorkflow.Name)

	// Test CreateWorkflowRun
	workflowRun := &entities.DBWorkflowRun{
		ID:           db.GenerateReadableID(),
		Input:        json.RawMessage(`{"test": "input"}`),
		WorkflowName: "test-workflow",
		Status:       entities.WorkflowStatusPending,
		ScheduledAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = repo.CreateWorkflowRun(ctx, workflowRun)
	require.NoError(t, err)

	// Test GetNextWorkflowRun
	nextRun, err := repo.GetNextWorkflowRun(ctx)
	require.NoError(t, err)
	require.NotNil(t, nextRun)

	assert.Equal(t, workflowRun.ID, nextRun.ID)
	assert.Equal(t, workflowRun.WorkflowName, nextRun.WorkflowName)
	assert.Equal(t, workflowRun.Status, nextRun.Status)
	assert.JSONEq(t, string(workflowRun.Input), string(nextRun.Input))

	// Test ChangeWorkflowRunStatus
	err = repo.ChangeWorkflowRunStatus(ctx, workflowRun.ID, entities.WorkflowStatusExecuting)
	require.NoError(t, err)

	// Verify status changed by getting the next workflow run (should be nil since it's executing)
	nextRunAfterUpdate, err := repo.GetNextWorkflowRun(ctx)
	require.NoError(t, err)
	require.Nil(t, nextRunAfterUpdate)
}
