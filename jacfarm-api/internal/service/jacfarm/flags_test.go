package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/internal/service/jacfarm/mocks"
	"JacFARM/pkg/rabbitmq_dto"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListFlags(t *testing.T) {
	suite := newTestSuite(t)

	testcases := []struct {
		name          string
		filter        *dto.ListFlagsFilter
		mock          func() *mocks.StorageMock
		expectedFlags []*models.FlagEnrich
		expectedCount int
		expectedErr   error
	}{
		{
			name: "ok",
			filter: &dto.ListFlagsFilter{
				Limit: 10,
				Page:  1,
			},
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetFlags(
					gomock.Any(),
					&dto.ListFlagsFilter{
						Limit: 10,
						Page:  1,
					},
				).Return(
					[]*models.FlagEnrich{
						{
							ID: int64(1),
						},
					},
					1,
					nil,
				)
				return suite.storageMock
			},
			expectedFlags: []*models.FlagEnrich{
				{
					ID: int64(1),
				},
			},
			expectedCount: 1,
		},
		{
			name: "internal error",
			filter: &dto.ListFlagsFilter{
				Limit: 10,
				Page:  1,
			},
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetFlags(
					gomock.Any(),
					&dto.ListFlagsFilter{
						Limit: 10,
						Page:  1,
					},
				).Return(nil, 0, testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageMock := tc.mock()
			suite.srv.db = storageMock
			flags, count, err := suite.srv.ListFlags(t.Context(), tc.filter)
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedFlags, flags)
			require.Equal(t, tc.expectedCount, count)
		})
	}
}

func TestPutFlag(t *testing.T) {
	suite := newTestSuite(t)

	testcases := []struct {
		name        string
		flag        string
		mock        func() *mocks.StorageMock
		expectedErr error
	}{
		{
			name: "ok",
			flag: "123",
			mock: func() *mocks.StorageMock {
				suite.queueMock.EXPECT().PublishFlag(
					gomock.AssignableToTypeOf(&rabbitmq_dto.Flag{}),
				).Return(nil)
				return suite.storageMock
			},
		},
		{
			name: "internal error",
			flag: "123",
			mock: func() *mocks.StorageMock {
				suite.queueMock.EXPECT().PublishFlag(
					gomock.AssignableToTypeOf(&rabbitmq_dto.Flag{}),
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
			err := suite.srv.PutFlag(t.Context(), tc.flag)
			require.ErrorIs(t, tc.expectedErr, err)
		})
	}
}

func TestGetFlagsCount(t *testing.T) {
	suite := newTestSuite(t)

	testcases := []struct {
		name        string
		count       int
		mock        func() *mocks.StorageMock
		expectedErr error
	}{
		{
			name:  "ok",
			count: 1,
			mock: func() *mocks.StorageMock {
				suite.queueMock.EXPECT().GetFlagsCount(gomock.Any()).Return(1, nil)
				return suite.storageMock
			},
		},
		{
			name:  "internal error",
			count: 0,
			mock: func() *mocks.StorageMock {
				suite.queueMock.EXPECT().GetFlagsCount(gomock.Any()).Return(0, testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageMock := tc.mock()
			suite.srv.db = storageMock
			count, err := suite.srv.GetFlagsCount(t.Context())
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.count, count)
		})
	}
}
