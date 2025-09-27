package txrepo_test

import (
	"context"
	"log"
	"testing"

	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/utils"
)

var testContainer *utils.PGTestContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start PostgreSQL container
	var err error
	testContainer, err = utils.GetPGTestContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to start test container: %v", err)
	}

	// Initialize database schema
	initiator := db.NewPGInitiator(testContainer.GetPool())
	err = initiator.Init(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Run tests
	m.Run()

	// Cleanup
	if err := testContainer.Close(ctx); err != nil {
		log.Printf("Failed to close test container: %v", err)
	}

}
