package server_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.eloylp.dev/go-serve/server"
)

func TestUploadTARGZHandlerAcceptsRelativeRoot(t *testing.T) {
	logBuffer := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(logBuffer)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/upload", sampleTARGZContentReader())

	deployPath := "teststuff/" + t.Name()

	defer os.RemoveAll("teststuff")
	req.Header.Add(DeployPathHeader, deployPath)
	req.Header.Add("Content-Type", server.ContentTypeTarGzip)
	server.UploadTARGZHandler(logger, ".").ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	// In case of test failure, just log the body and the logger output to get better context.
	data, err := io.ReadAll(rec.Result().Body)
	require.NoError(t, err)
	t.Log(string(data))
	t.Log(logBuffer.String())

}

func TestDownloadTARGZHandlerAcceptsRelativeRoot(t *testing.T) {
	logBuffer := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(logBuffer)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/download", nil)

	req.Header.Add(DownloadPathHeader, "handler_test.go")
	req.Header.Add("Accept", server.ContentTypeTarGzip)
	server.DownloadTARGZHandler(logger, ".").ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	// In case of test failure, just log the body and the logger output to get better context.
	data, err := io.ReadAll(rec.Result().Body)
	require.NoError(t, err)
	t.Log(string(data))
	t.Log(logBuffer.String())
}
