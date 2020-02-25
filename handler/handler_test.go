package handler_test

import (
	"github.com/eloylp/go-serve/handler"
	"github.com/eloylp/go-serve/logging/mock"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerHeader(t *testing.T) {

	rec := httptest.NewRecorder()
	middleware := handler.ServerHeader("v1.0.0")
	testHandler := handlerFixture()
	chain := middleware(testHandler)
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	chain.ServeHTTP(rec, request)
	assert.Equal(t, "go-serve v1.0.0", rec.Result().Header.Get("Server"),
		"Server header is not matching name version format")
	assertOriginalHandlerExecution(t, rec.Result().Body)
}

func handlerFixture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Handle wrote this"))
	}
}

func assertOriginalHandlerExecution(t *testing.T, body io.ReadCloser) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Handle wrote this", string(data), "Handler is not correctly executed")
}

func TestRequestLogger(t *testing.T) {

	rec := httptest.NewRecorder()
	fakeLogger := mock.NewFakeLogger()
	middleware := handler.RequestLogger(fakeLogger)
	testHandler := handlerFixture()

	chain := middleware(testHandler)
	request, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.RemoteAddr = "127.0.0.1"
	fakeLogger.On("Infof", "%s %s from client %s",
		request.Method, "/path", request.RemoteAddr).Return()

	chain.ServeHTTP(rec, request)
	fakeLogger.AssertExpectations(t)
	assertOriginalHandlerExecution(t, rec.Result().Body)
}
