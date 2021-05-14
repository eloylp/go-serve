//+build integration

package server_test

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertTARGZMD5Sums(t *testing.T, r io.Reader, expectedElems map[string]string) {
	gzipReader, err := gzip.NewReader(r)
	assert.NoError(t, err)
	tarReader := tar.NewReader(gzipReader)
	elems := map[string]string{}
	for {
		h, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		sum := md5.New()
		if !h.FileInfo().IsDir() {
			_, err = io.Copy(sum, tarReader)
			if err != nil {
				t.Fatal(err)
			}
			elems[h.Name] = fmt.Sprintf("%x", sum.Sum(nil))
			continue
		}
		elems[h.Name] = ""
	}
	assert.Equal(t, expectedElems, elems)
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
