//+build integration

package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTARGZUpload(t *testing.T) {
	BeforeEach(t)

	s, logBuff, _ := sut(t)

	defer s.Shutdown(context.Background())

	// Get a sample of compressed doc root. It will contain 2 images, tux.png and gnu.png.
	tarGZFile, err := os.Open(DocRootTARGZ)
	require.NoError(t, err)
	defer tarGZFile.Close()

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", tarGZFile)
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.Header.Add("GoServe-Deploy-Path", "/sub-root/test")

	// Send tar.gz to the upload endpoint
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	expectedSuccessMessage := "upload complete ! Bytes written: 533766"
	assert.Equal(t, expectedSuccessMessage, string(data))

	// Check that files are served correctly.
	tux := BodyFrom(t, HTTPAddressStatic+"/sub-root/test/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(tux), "got body: %s", tux)

	gnu := BodyFrom(t, HTTPAddressStatic+"/sub-root/test/gnu.png")
	assert.Equal(t, GnuTestFileMD5, md5From(gnu), "got body: %s", gnu)

	notes := BodyFrom(t, HTTPAddressStatic+"/sub-root/test/notes/notes.txt")
	assert.Equal(t, NotesTestFileMD5, md5From(notes), "got body: %s", notes)

	subNotes := BodyFrom(t, HTTPAddressStatic+"/sub-root/test/notes/subnotes/notes.txt")
	assert.Equal(t, SubNotesTestFileMD5, md5From(subNotes), "got body: %s", notes)
	s.Shutdown(context.Background()) // Force shutdown here in order to avoid data race with the logger buffer
	assert.Contains(t, logBuff.String(), expectedSuccessMessage)
}

func TestTARGZUploadCannotEscapeFromDocRoot(t *testing.T) {
	BeforeEach(t)

	s, logBuff, _ := sut(t)

	defer s.Shutdown(context.Background())

	tarGZFile, err := os.Open(DocRootTARGZ)
	require.NoError(t, err)
	defer tarGZFile.Close()

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", tarGZFile)
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.Header.Add("GoServe-Deploy-Path", "..")

	// Send tar.gz to the upload endpoint
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.NoError(t, err)
	logs := logBuff.String()
	assert.Contains(t, logs, "upload path violation try")
}

func TestUpload(t *testing.T) {
	BeforeEach(t)

	s, logBuff, _ := sut(t)

	defer s.Shutdown(context.Background())

	// Get a sample of compressed doc root. It will contain 2 images, tux.png and gnu.png.
	file, err := os.Open(DocRoot + "/notes/notes.txt")
	require.NoError(t, err)
	defer file.Close()

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", file)
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("GoServe-Deploy-Path", "/sub-root/notes.txt")

	// Send tar.gz to the upload endpoint
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	expectedSuccessMessage := "upload complete ! Bytes written: 20"
	assert.Equal(t, expectedSuccessMessage, string(data))

	// Check that file is served correctly.
	notes := BodyFrom(t, HTTPAddressStatic+"/sub-root/notes.txt")
	assert.Equal(t, NotesTestFileMD5, md5From(notes), "got body: %s", notes)

	s.Shutdown(context.Background()) // Force shutdown here in order to avoid data race with the logger buffer
	assert.Contains(t, logBuff.String(), expectedSuccessMessage)
}
