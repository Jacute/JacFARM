package jacfarm

import (
	"JacFARM/internal/models"
	"JacFARM/internal/service/jacfarm/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetStatuses(t *testing.T) {
	suite := newTestSuite(t)

	testcases := []struct {
		name             string
		mock             func() *mocks.StorageMock
		expectedStatuses []*models.Status
		expectedErr      error
	}{
		{
			name: "ok",
			expectedStatuses: []*models.Status{
				{
					ID:   1,
					Name: models.FlagStatusOld,
				},
			},
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetStatuses(gomock.Any()).Return([]*models.Status{
					{
						ID:   1,
						Name: models.FlagStatusOld,
					},
				}, nil)
				return suite.storageMock
			},
		},
		{
			name: "internal error",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetStatuses(gomock.Any()).Return(nil, testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageMock := tc.mock()
			suite.srv.db = storageMock
			statuses, err := suite.srv.GetStatuses(t.Context())
			require.Equal(t, tc.expectedStatuses, statuses)
			require.ErrorIs(t, tc.expectedErr, err)
		})
	}
}
