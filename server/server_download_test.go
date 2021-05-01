package server_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/eloylp/kit/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func TestTARGZDownload(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	testDocRoot := t.TempDir()
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(testDocRoot),
			config.WithLoggerOutput(logBuff),
			config.WithDownLoadEndpoint("/download"),
			config.WithLoggerLevel(logrus.DebugLevel.String()),
		),
	)
	assert.NoError(t, err)

	test.Copy(t, DocRoot, testDocRoot)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, HTTPAddress+"/download", nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", "application/tar+gzip")
	req.Header.Add("GoServe-Download-Path", "/notes")

	// Obtain the tar.gz with the required path
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
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
	logBuff := bytes.NewBuffer(nil)
	testDocRoot := t.TempDir()
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(testDocRoot),
			config.WithLoggerOutput(logBuff),
			config.WithDownLoadEndpoint("/download"),
			config.WithLoggerLevel(logrus.DebugLevel.String()),
		),
	)
	assert.NoError(t, err)

	test.Copy(t, DocRoot, testDocRoot)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, HTTPAddress+"/download", nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", "application/tar+gzip")
	req.Header.Add("GoServe-Download-Path", "/notes/notes.txt")

	// Obtain the tar.gz with the required path
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		AssertTARGZMD5Sums(t, resp.Body, map[string]string{
			"notes.txt": NotesTestFileMD5,
		})
	}
}

func TestTARGZDownloadCannotEscapeFromDocRoot(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(t.TempDir()),
			config.WithLoggerOutput(logBuff),
			config.WithDownLoadEndpoint("/download"),
			config.WithLoggerLevel(logrus.DebugLevel.String()),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	// Prepare request
	req, err := http.NewRequest(http.MethodGet, HTTPAddress+"/download", nil)
	assert.NoError(t, err)
	req.Header.Add("Accept", "application/tar+gzip")
	req.Header.Add("GoServe-Download-Path", "..")

	// Require tar.gz to the download endpoint
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	// We need to calculate one directory up to the doc root to check the message is correct.
	// This is due to the previous 	req.Header.Add("GoServe-Download-Path", "..") statement.
	docRootDirParts := filepath.SplitList(filepath.Dir(t.TempDir()))
	parentDocRoot := filepath.Join(docRootDirParts[0:]...)
	expectedMessage := fmt.Sprintf("the path you provided %s is not a suitable one", parentDocRoot)
	assert.Equal(t, expectedMessage, string(data))
	assert.Contains(t, logBuff.String(), expectedMessage)
	assert.Contains(t, logBuff.String(), "download path violation try")
}
