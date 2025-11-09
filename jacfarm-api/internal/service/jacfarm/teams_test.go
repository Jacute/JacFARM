package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/internal/service/jacfarm/mocks"
	storage_errors "JacFARM/internal/storage"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListShortTeams(t *testing.T) {
	suite := newTestSuite(t)

	testcases := []struct {
		name        string
		mock        func() *mocks.StorageMock
		expected    []*models.ShortTeam
		expectedErr error
	}{
		{
			name: "ok",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetShortTeams(gomock.Any()).
					Return([]*models.ShortTeam{{ID: 1, IP: net.ParseIP("10.10.1.1")}}, nil)
				return suite.storageMock
			},
			expected: []*models.ShortTeam{{ID: 1, IP: net.ParseIP("10.10.1.1")}},
		},
		{
			name: "db error",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetShortTeams(gomock.Any()).
					Return(nil, testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageMock := tc.mock()
			suite.srv.db = storageMock

			teams, err := suite.srv.ListShortTeams(t.Context())
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expected, teams)
		})
	}
}

func TestListTeams(t *testing.T) {
	suite := newTestSuite(t)
	filter := &dto.ListTeamsFilter{Limit: 5, Page: 1}

	testcases := []struct {
		name        string
		mock        func() *mocks.StorageMock
		expected    []*models.Team
		expectedCnt int
		expectedErr error
	}{
		{
			name: "ok",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetTeams(gomock.Any(), filter).
					Return([]*models.Team{{ID: 10, Name: "Team X"}}, 1, nil)
				return suite.storageMock
			},
			expected:    []*models.Team{{ID: 10, Name: "Team X"}},
			expectedCnt: 1,
		},
		{
			name: "db error",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().GetTeams(gomock.Any(), filter).
					Return(nil, 0, testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			suite.srv.db = tc.mock()

			teams, count, err := suite.srv.ListTeams(t.Context(), filter)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expected, teams)
			require.Equal(t, tc.expectedCnt, count)
		})
	}
}

func TestAddTeam(t *testing.T) {
	suite := newTestSuite(t)
	team := &models.Team{Name: "New Team"}

	testcases := []struct {
		name        string
		mock        func() *mocks.StorageMock
		expectedID  int64
		expectedErr error
	}{
		{
			name: "ok",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().AddTeam(gomock.Any(), team).
					Return(int64(42), nil)
				return suite.storageMock
			},
			expectedID: 42,
		},
		{
			name: "already exists",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().AddTeam(gomock.Any(), team).
					Return(int64(0), storage_errors.ErrTeamAlreadyExists)
				return suite.storageMock
			},
			expectedErr: storage_errors.ErrTeamAlreadyExists,
		},
		{
			name: "db error",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().AddTeam(gomock.Any(), team).
					Return(int64(0), testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			suite.srv.db = tc.mock()
			id, err := suite.srv.AddTeam(t.Context(), team)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedID, id)
		})
	}
}

func TestDeleteTeam(t *testing.T) {
	suite := newTestSuite(t)
	id := int64(123)

	testcases := []struct {
		name        string
		mock        func() *mocks.StorageMock
		expectedErr error
	}{
		{
			name: "ok",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().DeleteTeam(gomock.Any(), id).
					Return(nil)
				return suite.storageMock
			},
		},
		{
			name: "team not found",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().DeleteTeam(gomock.Any(), id).
					Return(storage_errors.ErrTeamNotFound)
				return suite.storageMock
			},
			expectedErr: storage_errors.ErrTeamNotFound,
		},
		{
			name: "db error",
			mock: func() *mocks.StorageMock {
				suite.storageMock.EXPECT().DeleteTeam(gomock.Any(), id).
					Return(testDbErr)
				return suite.storageMock
			},
			expectedErr: testDbErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			suite.srv.db = tc.mock()
			err := suite.srv.DeleteTeam(t.Context(), id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
