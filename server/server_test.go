package server_test

import (
	"bytes"
	"io"
	"os"

	"go.eloylp.dev/go-serve/server"
)

const (
	ListenAddress       = "localhost:9090"
	HTTPAddress         = "http://" + ListenAddress
	HTTPAddressStatic   = "http://" + ListenAddress + "/static"
	HTTPAddressUpload   = "http://" + ListenAddress + "/upload"
	HTTPAddressDownload = "http://" + ListenAddress + "/download"
	HTTPAddressStatus   = "http://" + ListenAddress + "/status"
	DocRoot             = "../tests/root"
	TuxTestFileMD5      = "a0e6e27f7e31fd0bd549ea936033bf28"
	GnuTestFileMD5      = "0073978283cb69d470ec2ea1b66f1988"
	NotesTestFileMD5    = "36d7e788e7a54109f5beb9ebe103da39"
	SubNotesTestFileMD5 = "0ff6da62cf7875cce432f7b955008953"
	DocRootTARGZ        = "../tests/doc-root.tar.gz"
	DeployPathHeader    = "GoServe-Deploy-Path"
	DownloadPathHeader  = "GoServe-Download-Path"
)

var (
	testUserCredentials = map[string]string{
		"user": "$2y$10$mAx10mlJ/UNbQJCgPp2oLe9n9jViYl9vlT0cYI3Nfop3P3bU1PDay", // Unencrypted value: user:password
	}
	sampleTARGZContent = func() []byte {
		file, err := os.ReadFile(DocRootTARGZ)
		if err != nil {
			panic(err)
		}
		return file
	}()
)

func sampleTARGZContentReader() io.Reader {
	return bytes.NewReader(sampleTARGZContent)
}

func init() {
	server.Name = "go-serve"
	server.Version = "v1.0.0"
	server.Build = "af09"
	server.BuildTime = "1988-01-21"
}
