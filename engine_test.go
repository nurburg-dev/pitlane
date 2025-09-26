package pitlane_test

import (
	"context"
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
