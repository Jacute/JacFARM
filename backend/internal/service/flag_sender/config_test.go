package flag_sender

import (
	jacfarm "JacFARM/internal/service/JacFARM"
	"JacFARM/internal/storage/sqlite"
	"context"
	"log/slog"
	"testing"

	"github.com/jacute/prettylogger"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	testCtx := context.Background()
	fakeLogger := slog.New(prettylogger.NewDiscardHandler())

	storage := sqlite.PrepareDBToTest(testCtx, t)

	flagSenderConfig := map[string]string{
		jacfarm.ConfigFlagSenderFlagTTL:       "1m",
		jacfarm.ConfigFlagSenderJuryFlagURL:   "http://localhost:8080",
		jacfarm.ConfigFlagSenderPlugin:        "test",
		jacfarm.ConfigFlagSenderSubmitLimit:   "1",
		jacfarm.ConfigFlagSenderSubmitPeriod:  "1m",
		jacfarm.ConfigFlagSenderSubmitTimeout: "1m",
		jacfarm.ConfigFlagSenderToken:         "test",
	}

	for k, v := range flagSenderConfig {
		err := storage.AddConfigParameter(testCtx, k, v)
		require.NoError(t, err)
	}

	fs, err := New(
		fakeLogger,
		storage,
		"test",
	)
	require.NotNil(t, fs)
	require.NoError(t, err)
	require.Equal(t, "test", fs.cfg.pluginDir)

	err = fs.loadConfig(testCtx, true)
	require.NoError(t, err)

	// TODO: add more tests
}
