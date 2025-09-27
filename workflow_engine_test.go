package pitlane_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/nurburg-dev/pitlane"
	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/utils"
	"github.com/stretchr/testify/require"
)

var (
	pgContainer *utils.PGTestContainer
	tables      = []string{"workflows", "workflow_runs"}
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	pgc, err := utils.GetPGTestContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	pgContainer = pgc
	defer func() {
		_ = pgc.Close(ctx)
	}()
	m.Run()
}

func TestEngineInit(t *testing.T) {
	ctx := context.Background()
	cfg := pitlane.NewDBConfig(
		pgContainer.GetHost(),
		pgContainer.GetPort(),
		pgContainer.GetUsername(),
		pgContainer.GetDatabase(),
		pgContainer.GetPassword(),
	)
	// Part 1: check if database tables are not created
	we0, err := pitlane.NewWorkflowEngine(ctx, pitlane.NewEngineConfig(cfg, false))
	require.NoError(t, err)
	require.NotNil(t, we0)
	for _, table := range tables {
		t1, err2 := db.TableExists(ctx, pgContainer.GetPool(), table)
		require.NoError(t, err2)
		require.False(t, t1)
	}

	// Part 2: check if tables are created
	we, err := pitlane.NewWorkflowEngine(ctx, pitlane.NewEngineConfig(cfg, true))
	require.NoError(t, err)
	require.NotNil(t, we)
	for _, table := range tables {
		t1, err2 := db.TableExists(ctx, pgContainer.GetPool(), table)
		require.NoError(t, err2)
		require.True(t, t1)
	}
}

func Activity1(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

func Activity2(_ context.Context, a, b string) (string, error) {
	return a + b, nil
}

func SampleWorkflow(_ context.Context, name string, count int) (string, error) {
	return fmt.Sprintf("Hello %s %d", name, count), nil
}

func TestInvokeWorkflow(t *testing.T) {
	ctx := context.Background()
	cfg := pitlane.NewDBConfig(
		pgContainer.GetHost(),
		pgContainer.GetPort(),
		pgContainer.GetUsername(),
		pgContainer.GetDatabase(),
		pgContainer.GetPassword(),
	)

	we, err := pitlane.NewWorkflowEngine(ctx, pitlane.NewEngineConfig(cfg, true))
	require.NoError(t, err)
	require.NotNil(t, we)

	// Register a workflow
	err = pitlane.RegisterWorkflow(SampleWorkflow)
	require.NoError(t, err)

	// Trigger the workflow
	workflowRunID, err := we.InvokeWorkflow(ctx, SampleWorkflow, "test", 42)
	require.NoError(t, err)
	require.NotEmpty(t, workflowRunID)

	// Verify workflow run was created in database
	query := `SELECT id, workflow_name, status, input FROM workflow_runs WHERE id = $1`
	row := pgContainer.GetPool().QueryRow(ctx, query, workflowRunID)

	var id, workflowName, status string
	var input []byte
	err = row.Scan(&id, &workflowName, &status, &input)
	require.NoError(t, err)
	require.Equal(t, workflowRunID, id)
	require.Equal(t, "github.com/nurburg-dev/pitlane_test.SampleWorkflow", workflowName)
	require.Equal(t, "pending", status)
	require.Contains(t, string(input), "test")
	require.Contains(t, string(input), "42")
}
