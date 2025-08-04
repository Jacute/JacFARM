package sqlite

import (
	"JacFARM/internal/models"
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const testDbPath = "./test.db"
const migrationsPath = "../../../migrations"

func TestUpdateStatusForOldFlags(t *testing.T) {
	testCtx := context.Background()

	storage, err := New(testDbPath)
	require.NotNil(t, storage)
	require.NoError(t, err)
	storage.ApplyMigrations(testCtx, testDbPath, migrationsPath)
	t.Cleanup(func() {
		storage.Close()
		os.Remove(testDbPath)
	})

	// create team
	teamID, err := storage.AddTeam(&models.Team{
		Name: "123",
		IP:   "10.10.1.2",
	})
	require.NoError(t, err)

	// create exploit
	exploitID := uuid.NewString()
	reqPath := "test"
	err = storage.CreateExploit(testCtx, &models.Exploit{
		ID:               exploitID,
		Name:             "test",
		Type:             "test",
		IsLocal:          true,
		ExecutablePath:   "test",
		RequirementsPath: &reqPath,
	})
	require.NoError(t, err)

	puttedFlag := &models.Flag{
		Value:             "test",
		Status:            models.FlagStatusPending,
		ExploitID:         exploitID,
		GetFrom:           teamID,
		MessageFromServer: "",
		CreatedAt:         time.Now().UTC().Unix(), // make it old
	}
	flagID, err := storage.PutFlag(testCtx, puttedFlag)
	require.NoError(t, err)

	flag, err := storage.GetFlagEnrichByValue(testCtx, puttedFlag.Value)
	require.NoError(t, err)
	require.Equal(t, flagID, flag.ID)

	time.Sleep(time.Second)

	count, err := storage.UpdateStatusForOldFlags(testCtx, time.Second)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
	flag, err = storage.GetFlagEnrichByValue(testCtx, puttedFlag.Value)
	require.NoError(t, err)
	require.Equal(t, models.FlagStatusOld, flag.Status)
}

func TestUpdateStatusForOldFlagsError(t *testing.T) {
	testCtx := context.Background()

	storage, err := New(testDbPath)
	require.NotNil(t, storage)
	require.NoError(t, err)
	storage.ApplyMigrations(testCtx, testDbPath, migrationsPath)
	t.Cleanup(func() {
		storage.Close()
		os.Remove(testDbPath)
	})

	// create team
	teamID, err := storage.AddTeam(&models.Team{
		Name: "123",
		IP:   "10.10.1.2",
	})
	require.NoError(t, err)

	// create exploit
	exploitID := uuid.NewString()
	reqPath := "test"
	err = storage.CreateExploit(testCtx, &models.Exploit{
		ID:               exploitID,
		Name:             "test",
		Type:             "test",
		IsLocal:          true,
		ExecutablePath:   "test",
		RequirementsPath: &reqPath,
	})
	require.NoError(t, err)

	puttedFlag := &models.Flag{
		Value:             "test",
		Status:            models.FlagStatusPending,
		ExploitID:         exploitID,
		GetFrom:           teamID,
		MessageFromServer: "",
		CreatedAt:         time.Now().UTC().Unix(), // make it old
	}
	flagID, err := storage.PutFlag(testCtx, puttedFlag)
	require.NoError(t, err)

	flag, err := storage.GetFlagEnrichByValue(testCtx, puttedFlag.Value)
	require.NoError(t, err)
	require.Equal(t, flagID, flag.ID)

	count, err := storage.UpdateStatusForOldFlags(testCtx, 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
	flag, err = storage.GetFlagEnrichByValue(testCtx, puttedFlag.Value)
	require.NoError(t, err)
	require.Equal(t, models.FlagStatusPending, flag.Status)
}
