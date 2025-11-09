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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
			expectedStatusCode: http.StatusOK,
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
			expectedStatusCode: http.StatusInternalServerError,
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
			expectedStatusCode: http.StatusBadRequest,
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
			expectedStatusCode: http.StatusBadRequest,
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
			expectedStatusCode: http.StatusBadRequest,
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
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			errorModel: dto.Error(dto.ErrPageIncorrectType.Error()),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			st := newTestSuite(t, tc.mock())
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
				require.Equal(t, tc.errorModel, resModel, tc.name)
				return
			}
			resModel := &dto.GetConfigResponse{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err)
			require.Equal(t, tc.resModel, resModel, tc.name)
		})
	}
}

func TestUpdateConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testcases := []struct {
		name               string
		id                 string
		body               []byte
		expectedStatusCode int
		mock               func() *mocks.MockService
		resModel           *dto.Response
		headers            map[string]string
	}{
		{
			name:               "ok",
			id:                 "1",
			body:               []byte(`{"value":"aboba"}`),
			expectedStatusCode: http.StatusOK,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().UpdateConfig(gomock.Any(), int64(1), "aboba").Return(nil)

				return serviceMock
			},
			resModel: dto.OK(),
			headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "invalid content-type",
			id:   "1",
			body: []byte(`{"value":"aboba"}`),
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			headers: map[string]string{
				"Content-Type": "text/html",
			},
			resModel:           dto.ErrInvalidContentType,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid id",
			id:   "das",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body:               []byte(`{"value":"aboba"}`),
			expectedStatusCode: http.StatusBadRequest,
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			resModel: dto.Error(dto.ErrIdIncorrectType.Error()),
		},
		{
			name: "config not found",
			id:   "1",
			body: []byte(`{"value":"aboba"}`),
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().UpdateConfig(gomock.Any(), int64(1), "aboba").Return(storage.ErrConfigParamNotFound)
				return serviceMock
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			resModel:           dto.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "service internal error",
			id:   "1",
			body: []byte(`{"value":"aboba"}`),
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				serviceMock.EXPECT().UpdateConfig(gomock.Any(), int64(1), "aboba").Return(fmt.Errorf("internal"))
				return serviceMock
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			resModel:           dto.ErrInternal,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "error decoding body",
			id:   "1",
			body: []byte(`{"value":"aboba`),
			mock: func() *mocks.MockService {
				serviceMock := mocks.NewMockService(ctrl)
				return serviceMock
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			resModel:           dto.ErrDecodingBody,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			st := newTestSuite(t, tc.mock())
			req := httptest.NewRequest(
				"PATCH",
				"/api/v1/config/"+tc.id,
				bytes.NewBuffer(tc.body),
			)
			for header, value := range tc.headers {
				req.Header.Add(header, value)
			}

			res, err := st.app.Test(req)
			defer res.Body.Close()
			require.NoError(t, err)

			data, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			resModel := &dto.Response{}
			err = json.Unmarshal(data, resModel)
			require.NoError(t, err, tc.name)
			require.Equal(t, tc.resModel, resModel, tc.name)

			require.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}
