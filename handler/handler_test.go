package handler_test

import (
	"github.com/eloylp/go-serve/handler"
	"github.com/eloylp/go-serve/logging/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersionHeader(t *testing.T) {

	rec := httptest.NewRecorder()
	middleware := handler.VersionHeader("v1.0.0")
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Handle wrote this"))
	})
	chain := middleware(testHandler)
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	chain.ServeHTTP(rec, request)
	assert.Equal(t, "go-serve v1.0.0", rec.Result().Header.Get("Server"),
		"Server header is not matching name version format")

	data, err := ioutil.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Handle wrote this", string(data), "Handler is not correctly executed")
}

func TestRequestLogger(t *testing.T) {

	rec := httptest.NewRecorder()
	fakeLogger := mock.NewFakeLogger()
	middleware := handler.RequestLogger(fakeLogger)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Handle wrote this"))
	})
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

	data, err := ioutil.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Handle wrote this", string(data))
}
