// +build racy

package server_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"go.eloylp.dev/kit/test"

	"go.eloylp.dev/go-serve/config"
)

var (
	tuxContent = func() []byte {
		open, err := os.Open(DocRoot + "/tux.png")
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(open)
		if err != nil {
			panic(err)
		}
		return data
	}()
	httpClient = func() *http.Client {
		tr := http.DefaultTransport.(*http.Transport).Clone()
		tr.MaxIdleConns = 500
		tr.MaxConnsPerHost = 500
		tr.MaxIdleConnsPerHost = 500
		return &http.Client{
			Transport: tr,
		}
	}()
)

func TestServerExecutionPaths(t *testing.T) {
	BeforeEach(t)
	server, _, docRoot := sut(t,
		config.WithReadAuthorizations(testUserCredentials),
		config.WithWriteAuthorizations(testUserCredentials),
	)
	defer server.Shutdown(context.Background())
	test.Copy(t, DocRoot, docRoot)

}

func uploadFile() {
	fileName := fmt.Sprintf("%s-tux.png", uuid.NewString())
	req, err := http.NewRequest(http.MethodGet, HTTPAddressUpload, bytes.NewReader(tuxContent))
	if err != nil {
		panic(err)
	}
	req.Header.Add(DeployPathHeader, "/"+fileName)
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Body.Close()
}

func downloadFile() {

}
