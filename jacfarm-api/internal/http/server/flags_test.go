package server

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/http/handlers/mocks"
	"JacFARM/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListFlags(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		queryParams        map[string][]string
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.GetFlagsResponse
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
				service.EXPECT().ListFlags(gomock.Any(), &dto.ListFlagsFilter{
					Limit: 10,
					Page:  1,
				}).Return(
					[]*models.FlagEnrich{
						{
							ID:    1,
							Value: "aboba",
						},
					}, 1, nil,
				)
				return service
			},
			resModel: &dto.GetFlagsResponse{
				Response: dto.OK(),
				Flags: []*models.FlagEnrich{
					{
						ID:    1,
						Value: "aboba",
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
				service.EXPECT().ListFlags(gomock.Any(), &dto.ListFlagsFilter{
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
			name: "team_id not number error",
			queryParams: map[string][]string{
				"limit":   []string{"10"},
				"page":    []string{"1"},
				"team_id": []string{"sda"},
			},
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			errorModel: dto.Error(dto.ErrTeamIdIncorrectType.Error()),
		},
		{
			name: "status_id not number error",
			queryParams: map[string][]string{
				"limit":     []string{"10"},
				"page":      []string{"1"},
				"status_id": []string{"sda"},
			},
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			errorModel: dto.Error(dto.ErrStatusIdIncorrectType.Error()),
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
				"/api/v1/flags?"+queryParamsToString(tc.queryParams),
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
			resModel := &dto.GetFlagsResponse{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
		})
	}
}

func TestPutFlag(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		body               []byte
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.Response
		headers            map[string]string
	}{
		{
			name:               "ok",
			body:               []byte(`{"flag":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="}`),
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().PutFlag(
					gomock.Any(),
					"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
				).Return(nil)
				return service
			},
			resModel: dto.OK(),
			headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:               "internal error",
			body:               []byte(`{"flag":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="}`),
			expectedStatusCode: http.StatusInternalServerError,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().PutFlag(
					gomock.Any(),
					"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
				).Return(errors.New("internal"))
				return service
			},
			resModel: dto.ErrInternal,
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
			resModel: dto.ErrInvalidContentType,
		},
		{
			name:               "decoding body error",
			body:               []byte(`{123`),
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				return mocks.NewMockService(ctrl)
			},
			resModel: dto.ErrDecodingBody,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ts := newTestSuite(t, tc.mock())

			req := httptest.NewRequestWithContext(
				t.Context(),
				"POST",
				"/api/v1/flags",
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

			resModel := &dto.Response{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
		})
	}
}

func TestPutFlagWithServiceToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name        string
		token       string
		shouldValid bool
	}{
		{
			name:        "valid token",
			token:       testApiKey,
			shouldValid: true,
		},
		{
			name:        "invalid token",
			token:       "invalid_token",
			shouldValid: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ts := newTestSuite(t, mocks.NewMockService(ctrl))

			req := httptest.NewRequestWithContext(
				t.Context(),
				"POST",
				"/api/v1/service/flags",
				nil,
			)
			req.Header.Add("Authorization", tc.token)

			res, err := ts.app.Test(req)
			require.NoError(t, err)

			if tc.shouldValid {
				require.NotEqual(t, http.StatusUnauthorized, res.StatusCode)
			} else {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			}
		})
	}
}

func TestGetStatuses(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.GetStatusesResponse
		errModel           *dto.Response
	}{
		{
			name:               "ok",
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().GetStatuses(
					gomock.Any(),
				).Return(
					[]*models.Status{
						{
							ID:   1,
							Name: models.FlagStatusOld,
						},
					},
					nil,
				)
				return service
			},
			resModel: &dto.GetStatusesResponse{
				Response: dto.OK(),
				Statuses: []*models.Status{
					{
						ID:   1,
						Name: models.FlagStatusOld,
					},
				},
			},
		},
		{
			name:               "internal error",
			expectedStatusCode: http.StatusInternalServerError,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().GetStatuses(
					gomock.Any(),
				).Return(nil, errors.New("internal"))
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
				"/api/v1/flags/statuses",
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
			resModel := &dto.GetStatusesResponse{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
		})
	}
}

func TestGetFlagsCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.GetFlagsCountResponse
		errModel           *dto.Response
	}{
		{
			name:               "ok",
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().GetFlagsCount(gomock.Any()).Return(
					10,
					nil,
				)
				return service
			},
			resModel: &dto.GetFlagsCountResponse{
				Response: dto.OK(),
				Count:    10,
			},
		},
		{
			name:               "internal error",
			expectedStatusCode: http.StatusInternalServerError,
			mock: func() *mocks.MockService {
				service := mocks.NewMockService(ctrl)
				service.EXPECT().GetFlagsCount(gomock.Any()).Return(0, errors.New("internal"))
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
				"/api/v1/flags/count",
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
			resModel := &dto.GetFlagsCountResponse{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
		})
	}
}
