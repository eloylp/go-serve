//+build integration

package server_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.eloylp.dev/kit/test"
)

func TestTARGZDownload(t *testing.T) {
	BeforeEach(t)

	s, _, testDocRoot := sut(t)

	test.Copy(t, DocRoot, testDocRoot)

	defer s.Shutdown(context.Background())

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, HTTPAddressDownload, nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/tar+gzip")
	req.Header.Add("GoServe-Download-Path", "/notes")

	// Obtain the tar.gz with the required path
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		AssertTARGZMD5Sums(t, resp.Body, map[string]string{
			".":                  "",
			"notes.txt":          NotesTestFileMD5,
			"subnotes":           "",
			"subnotes/notes.txt": SubNotesTestFileMD5,
		})
	}
}

func TestTARGZDownloadForSingleFile(t *testing.T) {
	BeforeEach(t)

	s, _, testDocRoot := sut(t)

	test.Copy(t, DocRoot, testDocRoot)

	defer s.Shutdown(context.Background())

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, HTTPAddressDownload, nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/tar+gzip")
	req.Header.Add("GoServe-Download-Path", "/notes/notes.txt")

	// Obtain the tar.gz with the required path
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		AssertTARGZMD5Sums(t, resp.Body, map[string]string{
			"notes.txt": NotesTestFileMD5,
		})
	}
}

func TestTARGZDownloadCannotEscapeFromDocRoot(t *testing.T) {
	BeforeEach(t)

	s, logBuff, _ := sut(t)

	defer s.Shutdown(context.Background())

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, HTTPAddressDownload, nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/tar+gzip")
	req.Header.Add("GoServe-Download-Path", "..")

	// Require tar.gz to the download endpoint
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	logs := logBuff.String()
	assert.Contains(t, logs, "download path violation try")
}
