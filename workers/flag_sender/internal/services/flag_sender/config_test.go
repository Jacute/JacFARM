package flag_sender

import (
	"context"
	mocks "flag_sender/internal/services/flag_sender/mocks"
	"flag_sender/pkg/common_config"
	"fmt"
	"log/slog"
	"path"
	"testing"

	"github.com/jacute/prettylogger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestConfig_HappyPath(t *testing.T) {
	testCtx := context.Background()
	fakeLogger := slog.New(prettylogger.NewDiscardHandler())

	ctl := gomock.NewController(t)
	storage := mocks.NewStorageMock(ctl)
	queue := mocks.NewQueueMock(ctl)
	storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderPlugin).Return("plugin", nil)
	storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderToken).Return("test", nil)
	storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderFlagTTL).Return("1m", nil)
	storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderJuryFlagURL).Return("http://localhost:8080", nil)
	storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderSubmitLimit).Return("1", nil)
	storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderSubmitPeriod).Return("1m", nil)
	storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderSubmitTimeout).Return("1m", nil)

	pluginPath := setupPlugin(t, TestPlugin2)
	dir, _ := path.Split(pluginPath)

	fs, err := New(
		fakeLogger,
		dir,
		queue,
		storage,
	)
	require.NoError(t, err)
	require.NotNil(t, fs)
}

func TestConfig_Errors(t *testing.T) {
	testCtx := context.Background()
	fakeLogger := slog.New(prettylogger.NewDiscardHandler())

	ctl := gomock.NewController(t)

	testcases := []struct {
		name string
		mock func() *mocks.StorageMock
	}{
		{
			name: "no plugin",
			mock: func() *mocks.StorageMock {
				storage := mocks.NewStorageMock(ctl)
				storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderPlugin).Return("", fmt.Errorf("some error"))
				return storage
			},
		},
		{
			name: "no token",
			mock: func() *mocks.StorageMock {
				storage := mocks.NewStorageMock(ctl)
				storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderPlugin).Return("test", nil)
				storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderJuryFlagURL).Return("test", nil)
				storage.EXPECT().GetConfigParameter(testCtx, common_config.ConfigFlagSenderToken).Return("", fmt.Errorf("some error"))
				return storage
			},
		},
	}

	queue := mocks.NewQueueMock(ctl)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storage := tc.mock()
			fs, err := New(
				fakeLogger,
				"test",
				queue,
				storage,
			)
			require.Nil(t, fs)
			require.Error(t, err)
		})
	}
}
