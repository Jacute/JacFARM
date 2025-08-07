package sqlite

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDbPath = "./test.db"
const migrationsPath = "../../../migrations"

func PrepareDBToTest(ctx context.Context, t *testing.T) *Storage {
	storage, err := New(testDbPath)
	require.NotNil(t, storage)
	require.NoError(t, err)
	storage.ApplyMigrations(ctx, testDbPath, migrationsPath)
	t.Cleanup(func() {
		storage.Close()
		os.Remove(testDbPath)
	})

	return storage
}
