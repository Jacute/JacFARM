package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/internal/service/jacfarm/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetConfig(t *testing.T) {
	suite := newTestSuite(t)

	testcases := []struct {
		name           string
		filter         *dto.GetConfigFilter
		mock           func() *mocks.StorageMock
		expectedConfig []*models.Config
		expectedCount  int
		expectedErr    error
	}{
		{
			name: "ok",
			filter: &dto.GetConfigFilter{
				Limit: 10,
				Page:  1,
			},
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetConfig(
					gomock.Any(),
					&dto.GetConfigFilter{
						Limit: 10,
						Page:  1,
					},
				).Return([]*models.Config{
					{
						ID:    1,
						Name:  "123",
						Value: "456",
					},
				}, 1, nil)
				return suite.storageMock
			},
			expectedConfig: []*models.Config{
				{
					ID:    1,
					Name:  "123",
					Value: "456",
				},
			},
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name: "internal error",
			filter: &dto.GetConfigFilter{
				Limit: 10,
				Page:  1,
			},
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetConfig(
					gomock.Any(),
					&dto.GetConfigFilter{
						Limit: 10,
						Page:  1,
					},
				).Return(nil, 0, testDbErr)
				return suite.storageMock
			},
			expectedCount: 0,
			expectedErr:   testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageMock := tc.mock()
			suite.srv.db = storageMock
			exploits, count, err := suite.srv.GetConfig(t.Context(), tc.filter)
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedConfig, exploits)
			require.Equal(t, tc.expectedCount, count)
		})
	}
}

func TestUpdateConfig(t *testing.T) {
	suite := newTestSuite(t)

	testcases := []struct {
		name        string
		id          int64
		value       string
		mock        func() *mocks.StorageMock
		expectedErr error
	}{
		{
			name:  "ok",
			id:    1,
			value: "123",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().UpdateConfigRow(
					gomock.Any(),
					int64(1),
					"123",
				).Return(nil)
				return suite.storageMock
			},
			expectedErr: nil,
		},
		{
			name:  "internal error",
			id:    1,
			value: "123",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().UpdateConfigRow(
					gomock.Any(),
					int64(1),
					"123",
				).Return(testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageMock := tc.mock()
			suite.srv.db = storageMock
			err := suite.srv.UpdateConfig(t.Context(), tc.id, tc.value)
			require.ErrorIs(t, tc.expectedErr, err)
		})
	}
}
