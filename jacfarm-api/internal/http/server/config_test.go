package server

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/http/handlers/mocks"
	"JacFARM/internal/models"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	testcases := []struct {
		name               string
		queryParams        map[string][]string
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.GetConfigResponse
		errorModel         *dto.Response
	}{
		{
			name: "ok",
			queryParams: map[string][]string{
				"limit": []string{"10"},
				"page":  []string{"1"},
			},
			expectedStatusCode: 200,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return([]*models.Config{{
					ID:    1,
					Name:  "aboba",
					Value: "aboba",
				}}, 1, nil)

				return serviceMock
			},
			resModel: &dto.GetConfigResponse{
				Response: dto.OK(),
				Config: []*models.Config{{
					ID:    1,
					Name:  "aboba",
					Value: "aboba",
				}},
				Count: 1,
			},
		},
		{
			name: "service error",
			queryParams: map[string][]string{
				"limit": []string{"10"},
				"page":  []string{"1"},
			},
			expectedStatusCode: 500,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(nil, 0, errors.New("internal error"))

				return serviceMock
			},
			errorModel: dto.ErrInternal,
		},
		{
			name: "negative page error",
			queryParams: map[string][]string{
				"limit": []string{"10"},
				"page":  []string{"-1"},
			},
			expectedStatusCode: 400,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			errorModel: dto.Error(dto.ErrPageNegative.Error()),
		},
		{
			name: "negative limit error",
			queryParams: map[string][]string{
				"limit": []string{"-10"},
				"page":  []string{"1"},
			},
			expectedStatusCode: 400,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			errorModel: dto.Error(dto.ErrLimitNegative.Error()),
		},
		{
			name: "incorrect type limit error",
			queryParams: map[string][]string{
				"limit": []string{"-das10"},
				"page":  []string{"1"},
			},
			expectedStatusCode: 400,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			errorModel: dto.Error(dto.ErrLimitIncorrectType.Error()),
		},
		{
			name: "incorrect type page error",
			queryParams: map[string][]string{
				"limit": []string{"10"},
				"page":  []string{"das"},
			},
			expectedStatusCode: 400,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			errorModel: dto.Error(dto.ErrPageIncorrectType.Error()),
		},
	}

	for _, tc := range testcases {
		st := NewTestSuite(t, tc.mock())
		req := httptest.NewRequest(
			"GET",
			"/api/v1/config?"+queryParamsToString(tc.queryParams),
			nil,
		)

		res, err := st.app.Test(req)
		defer res.Body.Close()
		require.NoError(t, err)
		require.Equal(t, tc.expectedStatusCode, res.StatusCode)

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		if tc.errorModel != nil {
			resModel := &dto.Response{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.errorModel, resModel)
			continue
		}
		resModel := &dto.GetConfigResponse{}
		err = json.Unmarshal(data, resModel)
		require.NoError(t, err)
		require.Equal(t, tc.resModel, resModel)
	}
}
