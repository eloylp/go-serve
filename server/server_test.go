package server_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/eloylp/kit/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

const (
	ListenAddress    = "localhost:9090"
	HTTPAddress      = "http://" + ListenAddress
	DocRoot          = "../tests/root"
	TuxTestFileMD5   = "a0e6e27f7e31fd0bd549ea936033bf28"
	GnuTestFileMD5   = "0073978283cb69d470ec2ea1b66f1988"
	NotesTestFileMD5 = "36d7e788e7a54109f5beb9ebe103da39"
	DocRootTARGZ     = "../tests/doc-root.tar.gz"
)

func init() {
	server.Name = "programName"
	server.Version = "v1.0.0"
	server.Build = "af09"
	server.BuildTime = "1988-01-21"
}

func TestServingContentDefaults(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(DocRoot),
			config.WithDocRootPrefix("/"),
			config.WithLoggerOutput(logBuff),
		),
	)
	assert.NoError(t, err)
	go s.ListenAndServe()

	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	data := BodyFrom(t, HTTPAddress+"/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
	logs := logBuff.String()
	AssertNoProblemsInLogs(t, logs)
	AssertStartupLogs(t, logs)
	AssertShutdownLogs(t, logs)
}

func AssertNoProblemsInLogs(t *testing.T, logs string) {
	assert.NotContains(t, logs, "level=warning")
	assert.NotContains(t, logs, "level=error")
}

func AssertStartupLogs(t *testing.T, logs string) {
	assert.Contains(t, logs, "programName v1.0.0 af09 1988-01-21")
	assert.Contains(t, logs, "starting to serve")
}

func AssertShutdownLogs(t *testing.T, logs string) {
	assert.Contains(t, logs, "started gracefully shutdown of server ...")
	assert.Contains(t, logs, "server is now shutdown !")
}

func TestServingContentAlternatePath(t *testing.T) {
	logBuff := bytes.NewBuffer(nil)
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithDocRoot(DocRoot),
			config.WithDocRootPrefix("/alternate"),
			config.WithLoggerOutput(logBuff),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	data := BodyFrom(t, HTTPAddress+"/alternate/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(data), "got body: %s", data)
	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
	logs := logBuff.String()
	AssertNoProblemsInLogs(t, logs)
	AssertStartupLogs(t, logs)
	AssertShutdownLogs(t, logs)
}

func TestSeverIdentity(t *testing.T) {
	s, err := server.New(
		config.ForOptions(
			config.WithListenAddr(ListenAddress),
			config.WithLoggerOutput(ioutil.Discard),
		),
	)
	assert.NoError(t, err)

	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	test.WaitTCPService(t, ListenAddress, time.Millisecond, time.Second)

	resp, err := http.Get(HTTPAddress)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "go-serve v1.0.0", resp.Header.Get("server"))
}

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
	expectedSuccessMessage := "upload of tar.gz complete ! Bytes written: 533742"
	assert.Equal(t, expectedSuccessMessage, string(data))

	// Check that files are served correctly.
	tux := BodyFrom(t, HTTPAddress+"/sub-root/test/tux.png")
	assert.Equal(t, TuxTestFileMD5, md5From(tux), "got body: %s", tux)

	gnu := BodyFrom(t, HTTPAddress+"/sub-root/test/gnu.png")
	assert.Equal(t, GnuTestFileMD5, md5From(gnu), "got body: %s", gnu)

	notes := BodyFrom(t, HTTPAddress+"/sub-root/test/notes/notes.txt")
	assert.Equal(t, NotesTestFileMD5, md5From(notes), "got body: %s", notes)

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
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	// We need to calculate one directory up to the doc root to check the message is correct.
	// This is due to the previous 	req.Header.Add("GoServe-Deploy-Path", "..") statement.
	docRootDirParts := filepath.SplitList(filepath.Dir(t.TempDir()))
	parentDocRoot := filepath.Join(docRootDirParts[0:]...)
	expectedMessage := fmt.Sprintf("incorrect upload path: the path you provided %s/gnu.png is not a suitable one", parentDocRoot)
	assert.Equal(t, expectedMessage, string(data))
}
