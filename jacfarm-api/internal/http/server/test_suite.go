package server

import (
	"JacFARM/internal/config"
	"JacFARM/internal/http/handlers"
	"bytes"
	"log/slog"
	"mime/multipart"
	"net/url"
	"testing"
	"time"

	"github.com/jacute/prettylogger"
	"github.com/stretchr/testify/require"
)

const testApiKey = "test"

type testSuite struct {
	app *HTTPServer
}

type formFile struct {
	Name    string
	Content []byte
}

func newTestSuite(t *testing.T, s handlers.Service) *testSuite {
	h := handlers.New(s)
	log := slog.New(prettylogger.NewDiscardHandler())
	cfg := &config.HTTPConfig{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
		CORS: &config.CORSConfig{
			AllowedOrigins: []string{"http://test"},
		},
	}
	app := New(
		log,
		cfg,
		testApiKey,
		h,
	)
	go app.Start()
	t.Cleanup(func() {
		app.Stop()
	})

	return &testSuite{
		app: app,
	}
}

func queryParamsToString(queryParams map[string][]string) string {
	q := url.Values{}
	for k, vs := range queryParams {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	return q.Encode()
}

func createMultipartBody(
	t *testing.T,
	fields map[string]string,
	fileFields map[string]*formFile, // fileName -> content
) (body *bytes.Buffer, contentType string) {
	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// add simple fields
	for k, v := range fields {
		err := writer.WriteField(k, v)
		require.NoError(t, err)
	}

	// add file fields
	for fileField, formfile := range fileFields {
		fw, err := writer.CreateFormFile(fileField, formfile.Name)
		require.NoError(t, err)

		_, err = fw.Write(formfile.Content)
		require.NoError(t, err)
	}

	// close writer to finalize boundary
	err := writer.Close()
	require.NoError(t, err)

	return body, writer.FormDataContentType()
}
