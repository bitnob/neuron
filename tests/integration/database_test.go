package integration

import (
	"context"
	"testing"
)

func testDatabaseOperations(t *testing.T, suite *IntegrationSuite) {
	ctx := context.Background()

	// Test transaction handling
	t.Run("Transaction", func(t *testing.T) {
		tx, err := suite.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Test rollback
		if err := tx.Rollback(); err != nil {
			t.Errorf("Failed to rollback transaction: %v", err)
		}

		// Test commit with new transaction
		tx, err = suite.db.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Errorf("Failed to commit transaction: %v", err)
		}
	})
}
