package server

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/http/handlers/mocks"
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListShortTeams(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.ListShortTeamsResponse
		errModel           *dto.Response
	}{
		{
			name:               "ok",
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().ListShortTeams(gomock.Any()).Return(
					[]*models.ShortTeam{
						{
							ID: 1,
							IP: net.ParseIP("10.10.5.2"),
						},
					},
					nil,
				)
				return service
			},
			resModel: &dto.ListShortTeamsResponse{
				Response: dto.OK(),
				Teams: []*models.ShortTeam{
					{
						ID: 1,
						IP: net.ParseIP("10.10.5.2"),
					},
				},
			},
		},
		{
			name:               "internal error",
			expectedStatusCode: http.StatusInternalServerError,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().ListShortTeams(gomock.Any()).Return(nil, errors.New("internal"))
				return service
			},
			errModel: dto.ErrInternal,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ts := newTestSuite(t, tc.mock())

			req := httptest.NewRequestWithContext(
				t.Context(),
				"GET",
				"/api/v1/teams/short",
				nil,
			)

			res, err := ts.app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if tc.errModel != nil {
				resModel := &dto.Response{}
				err = json.Unmarshal(data, resModel)
				require.NoError(t, err)
				require.Equal(t, tc.errModel, resModel, tc.name)
				return
			}
			resModel := &dto.ListShortTeamsResponse{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
		})
	}
}

func TestListTeams(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		queryParams        map[string][]string
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.ListTeamsResponse
		errorModel         *dto.Response
	}{
		{
			name: "ok",
			queryParams: map[string][]string{
				"limit": []string{"10"},
				"page":  []string{"1"},
			},
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().ListTeams(gomock.Any(), &dto.ListTeamsFilter{
					Limit: 10,
					Page:  1,
				}).Return(
					[]*models.Team{
						{
							ID:   1,
							Name: "aboba",
							IP:   net.ParseIP("10.10.5.2"),
						},
					}, 1, nil,
				)
				return service
			},
			resModel: &dto.ListTeamsResponse{
				Response: dto.OK(),
				Teams: []*models.Team{
					{
						ID:   1,
						Name: "aboba",
						IP:   net.ParseIP("10.10.5.2"),
					},
				},
				Count: 1,
			},
		},
		{
			name: "internal error",
			queryParams: map[string][]string{
				"limit": []string{"10"},
				"page":  []string{"1"},
			},
			expectedStatusCode: http.StatusInternalServerError,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().ListTeams(gomock.Any(), &dto.ListTeamsFilter{
					Limit: 10,
					Page:  1,
				}).Return(
					nil, 0, fmt.Errorf("internal"),
				)
				return service
			},
			errorModel: dto.ErrInternal,
		},
		{
			name: "page not number error",
			queryParams: map[string][]string{
				"limit": []string{"10"},
				"page":  []string{"dasdsa"},
			},
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			errorModel: dto.Error(dto.ErrPageIncorrectType.Error()),
		},
		{
			name: "limit not number error",
			queryParams: map[string][]string{
				"limit": []string{"das"},
				"page":  []string{"1"},
			},
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			errorModel: dto.Error(dto.ErrLimitIncorrectType.Error()),
		},
		{
			name: "limit negative number error",
			queryParams: map[string][]string{
				"limit": []string{"-5"},
				"page":  []string{"1"},
			},
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			errorModel: dto.Error(dto.ErrLimitNegative.Error()),
		},
		{
			name: "page negative number error",
			queryParams: map[string][]string{
				"limit": []string{"5"},
				"page":  []string{"-1"},
			},
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			errorModel: dto.Error(dto.ErrPageNegative.Error()),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ts := newTestSuite(t, tc.mock())

			req := httptest.NewRequestWithContext(
				t.Context(),
				"GET",
				"/api/v1/teams?"+queryParamsToString(tc.queryParams),
				nil,
			)

			res, err := ts.app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if tc.errorModel != nil {
				resModel := &dto.Response{}
				err = json.Unmarshal(data, resModel)
				require.NoError(t, err)
				require.Equal(t, tc.errorModel, resModel, tc.name)
				return
			}
			resModel := &dto.ListTeamsResponse{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
		})
	}
}

func TestAddTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		body               []byte
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.AddTeamResponse
		errModel           *dto.Response
		headers            map[string]string
	}{
		{
			name:               "ok",
			body:               []byte(`{"name":"aboba","ip":"10.10.5.2"}`),
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().AddTeam(
					gomock.Any(),
					&models.Team{
						Name: "aboba",
						IP:   net.ParseIP("10.10.5.2"),
					},
				).Return(int64(1), nil)
				return service
			},
			resModel: &dto.AddTeamResponse{
				Response: dto.OK(),
				ID:       1,
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:               "internal error",
			body:               []byte(`{"name":"aboba","ip":"10.10.5.2"}`),
			expectedStatusCode: http.StatusInternalServerError,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().AddTeam(
					gomock.Any(),
					&models.Team{
						Name: "aboba",
						IP:   net.ParseIP("10.10.5.2"),
					},
				).Return(int64(0), errors.New("internal error"))
				return service
			},
			errModel: dto.ErrInternal,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:               "incorrect content-type",
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			headers: map[string]string{
				"Content-Type": "application/xml",
			},
			errModel: dto.ErrInvalidContentType,
		},
		{
			name:               "decoding body error",
			body:               []byte(`{123`),
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			errModel: dto.ErrDecodingBody,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:               "invalid ip",
			body:               []byte(`{"name":"aboba","ip":"10.10.277.2"}`),
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			errModel: dto.Error("field IP should be correct ip address"),
		},
		{
			name:               "team with this ip already exists",
			body:               []byte(`{"name":"aboba","ip":"10.10.5.2"}`),
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().AddTeam(
					gomock.Any(),
					&models.Team{
						Name: "aboba",
						IP:   net.ParseIP("10.10.5.2"),
					},
				).Return(int64(0), storage.ErrTeamAlreadyExists)
				return service
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			errModel: dto.Error(storage.ErrTeamAlreadyExists.Error()),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ts := newTestSuite(t, tc.mock())

			req := httptest.NewRequestWithContext(
				t.Context(),
				"POST",
				"/api/v1/teams",
				bytes.NewBuffer(tc.body),
			)
			for k, v := range tc.headers {
				req.Header.Add(k, v)
			}

			res, err := ts.app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if tc.errModel != nil {
				resModel := &dto.Response{}
				err = json.Unmarshal(data, resModel)
				require.NoError(t, err)
				require.Equal(t, tc.errModel, resModel)
				return
			}

			resModel := &dto.AddTeamResponse{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel)
		})
	}
}

func TestDeleteTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		id                 string
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.Response
	}{
		{
			name:               "ok",
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().DeleteTeam(
					gomock.Any(),
					int64(1),
				).Return(nil)

				return serviceMock
			},
			resModel: dto.OK(),
		},
		{
			name:               "internal error",
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().DeleteTeam(
					gomock.Any(),
					int64(1),
				).Return(fmt.Errorf("internal"))

				return serviceMock
			},
			resModel: dto.ErrInternal,
		},
		{
			name:               "team not found",
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().DeleteTeam(
					gomock.Any(),
					int64(1),
				).Return(storage.ErrTeamNotFound)

				return serviceMock
			},
			resModel: dto.Error(storage.ErrTeamNotFound.Error()),
		},
		{
			name:               "negative id",
			id:                 "-1",
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			resModel: dto.Error(dto.ErrIdShouldBePos.Error()),
		},
		{
			name:               "str id",
			id:                 "dsadsa",
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			resModel: dto.Error(dto.ErrIdIncorrectType.Error()),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			st := newTestSuite(t, tc.mock())
			req := httptest.NewRequest(
				"DELETE",
				fmt.Sprintf("/api/v1/teams/%s", tc.id),
				nil,
			)

			res, err := st.app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			resModel := &dto.Response{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
			defer res.Body.Close()
		})
	}
}
