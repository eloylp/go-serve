package server_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/eloylp/kit/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func TestTARGZUpload(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(t.TempDir()),
			config.WithLoggerOutput(logBuff),
			config.WithUploadEndpoint("/upload"),
			config.WithLoggerLevel(logrus.DebugLevel.String()),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	// Get a sample of compressed doc root. It will contain 2 images, tux.png and gnu.png.
	tarGZFile, err := os.Open(DocRootTARGZ)
	assert.NoError(t, err)
	defer tarGZFile.Close()

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", tarGZFile)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.Header.Add("GoServe-Deploy-Path", "/sub-root/test")

	// Send tar.gz to the upload endpoint
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedSuccessMessage := "upload of tar.gz complete ! Bytes written: 533766"
	assert.Equal(t, expectedSuccessMessage, string(data))

	// Check that files are served correctly.
	tux := BodyFrom(t, HTTPAddress+"/sub-root/test/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(tux), "got body: %s", tux)

	gnu := BodyFrom(t, HTTPAddress+"/sub-root/test/gnu.png")
	assert.Equal(t, GnuTestFileMD5, md5From(gnu), "got body: %s", gnu)

	notes := BodyFrom(t, HTTPAddress+"/sub-root/test/notes/notes.txt")
	assert.Equal(t, NotesTestFileMD5, md5From(notes), "got body: %s", notes)

	subNotes := BodyFrom(t, HTTPAddress+"/sub-root/test/notes/subnotes/notes.txt")
	assert.Equal(t, SubNotesTestFileMD5, md5From(subNotes), "got body: %s", notes)
	s.Shutdown(context.Background()) // Force shutdown here in order to avoid data race with the logger buffer
	assert.Contains(t, logBuff.String(), expectedSuccessMessage)
}

func TestTARGZUploadCannotEscapeFromDocRoot(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(t.TempDir()),
			config.WithLoggerOutput(logBuff),
			config.WithUploadEndpoint("/upload"),
			config.WithLoggerLevel(logrus.DebugLevel.String()),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	tarGZFile, err := os.Open(DocRootTARGZ)
	assert.NoError(t, err)
	defer tarGZFile.Close()

	// Prepare request
	req, err := http.NewRequest(http.MethodPost, HTTPAddress+"/upload", tarGZFile)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/tar+gzip")
	req.Header.Add("GoServe-Deploy-Path", "..")

	// Send tar.gz to the upload endpoint
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NoError(t, err)
	logs := logBuff.String()
	assert.Contains(t, logs, "upload path violation try")
}
