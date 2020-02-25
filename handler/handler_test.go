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

const handlerFixtureBody = "Handle wrote this"

func TestServerHeader(t *testing.T) {

	rec := httptest.NewRecorder()
	middleware := handler.ServerHeader("v1.0.0")
	h := handlerFixture(t)
	chain := middleware(h)
	request := newTestRequest(t, "GET", "/", nil)

	chain.ServeHTTP(rec, request)

	assert.Equal(t, "go-serve v1.0.0", rec.Result().Header.Get("Server"),
		"Server header is not matching name version format")
	assertHandlerFixtureExecution(t, rec.Result().Body)
}

func handlerFixture(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(handlerFixtureBody)); err != nil {
			t.Fatal(err)
		}
	}
}

func assertHandlerFixtureExecution(t *testing.T, body io.ReadCloser) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, handlerFixtureBody, string(data), "Handler is not correctly executed")
}

func newTestRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func TestRequestLogger(t *testing.T) {

	rec := httptest.NewRecorder()
	logger := mock.NewFakeLogger()
	middleware := handler.RequestLogger(logger)
	h := handlerFixture(t)

	chain := middleware(h)
	request := newTestRequest(t, "GET", "/path", nil)
	request.RemoteAddr = "127.0.0.1"
	logger.On("Infof", "%s %s from client %s",
		request.Method, "/path", request.RemoteAddr).Return()

	chain.ServeHTTP(rec, request)

	logger.AssertExpectations(t)
	assertHandlerFixtureExecution(t, rec.Result().Body)
}
