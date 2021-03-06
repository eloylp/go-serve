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
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.eloylp.dev/kit/exec"
	"go.eloylp.dev/kit/test"

	"go.eloylp.dev/go-serve/config"
	"go.eloylp.dev/go-serve/server"
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
	uploadCount        int32
	uploadTotalTime    float64
	uploadTotalTimeL   sync.Mutex
	downloadCount      int32
	downloadTotalTime  float64
	downloadTotalTimeL sync.Mutex
)

func TestServerExecutionPaths(t *testing.T) {
	BeforeEach(t)
	s, _, docRoot := sut(t,
		config.WithReadAuthorizations(testUserCredentials),
		config.WithWriteAuthorizations(testUserCredentials),
	)
	defer s.Shutdown(context.Background())
	ctx, cancl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancl()
	test.Copy(t, DocRoot, docRoot)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-ctx.Done()
		wg.Done()
	}()
	go exec.Parallelize(ctx, wg, 10, uploadFile)
	go exec.Parallelize(ctx, wg, 10, downloadFile)
	wg.Wait()

	t.Log("Racy test stats >>")
	t.Logf("Totals: uploaded %v . downloaded %v \n",
		uploadCount, downloadCount)
	t.Logf("Average time in seconds: uploaded %v . downloaded %v \n",
		uploadTotalTime/float64(uploadCount), downloadTotalTime/float64(downloadCount))
}

func uploadFile() {
	fileName := fmt.Sprintf("%v-tux.png", atomic.LoadInt32(&uploadCount))
	req, err := http.NewRequest(http.MethodPost, HTTPAddressUpload, bytes.NewReader(tuxContent))
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth("user", "password")
	req.Header.Add(DeployPathHeader, "/"+fileName)
	req.Header.Add("Content-Type", server.ContentTypeFile)
	now := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(now)
	uploadTotalTimeL.Lock()
	uploadTotalTime += elapsed.Seconds()
	uploadTotalTimeL.Unlock()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("expected 200 got %v", resp.StatusCode))
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	atomic.AddInt32(&uploadCount, 1)
}

func downloadFile() {
	req, err := http.NewRequest(http.MethodGet, HTTPAddressStatic+"/gnu.png", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth("user", "password")
	now := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(now)
	downloadTotalTimeL.Lock()
	downloadTotalTime += elapsed.Seconds()
	downloadTotalTimeL.Unlock()
	resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("expected 200 got %v", resp.StatusCode))
	}
	atomic.AddInt32(&downloadCount, 1)
}
