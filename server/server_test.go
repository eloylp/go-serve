package server_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/server"
)

const (
	ListenAddress       = "localhost:9090"
	HTTPAddress         = "http://" + ListenAddress
	DocRoot             = "../tests/root"
	TuxTestFileMD5      = "a0e6e27f7e31fd0bd549ea936033bf28"
	GnuTestFileMD5      = "0073978283cb69d470ec2ea1b66f1988"
	NotesTestFileMD5    = "36d7e788e7a54109f5beb9ebe103da39"
	SubNotesTestFileMD5 = "0ff6da62cf7875cce432f7b955008953"
	DocRootTARGZ        = "../tests/doc-root.tar.gz"
)

var sampleTARGZContent = func() []byte {
	file, err := os.ReadFile(DocRootTARGZ)
	if err != nil {
		panic(err)
	}
	return file
}()

func init() {
	server.Name = "programName"
	server.Version = "v1.0.0"
	server.Build = "af09"
	server.BuildTime = "1988-01-21"
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
