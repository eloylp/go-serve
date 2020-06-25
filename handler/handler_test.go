//nolint:bodyclose
package handler_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/handler"
	"github.com/eloylp/go-serve/logging/mock"
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

func assertHandlerFixtureExecution(t *testing.T, body io.Reader) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, handlerFixtureBody, string(data), "Handler is not correctly executed")
}

func assertBodyContent(t *testing.T, expected string, body io.ReadCloser) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, string(data))
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

const htpasswdTestFile = "./.htpasswd-test"
const realm = "secured-server"

func TestAuthChecker_Valid(t *testing.T) {
	rec := httptest.NewRecorder()
	auth := handler.AuthChecker(realm, htpasswdTestFile)
	h := handlerFixture(t)
	chain := auth(h)
	request := httptest.NewRequest("GET", "/path", nil)
	request.SetBasicAuth("user", "abc1234")
	chain.ServeHTTP(rec, request)
	assert.Equal(t, rec.Result().StatusCode, http.StatusOK)
	assertHandlerFixtureExecution(t, rec.Result().Body)
}

func TestAuthChecker_ValidWithEmptyAuthPath(t *testing.T) {
	rec := httptest.NewRecorder()
	// No valid htpasswd file config file, so handler must fail
	auth := handler.AuthChecker(realm, "")
	h := handlerFixture(t)
	chain := auth(h)
	request := httptest.NewRequest("GET", "/path", nil)
	request.SetBasicAuth("user", "password")
	chain.ServeHTTP(rec, request)
	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	assertBodyContent(t, "Bad auth file", rec.Result().Body)
}

func TestAuthChecker_NotValidAuth(t *testing.T) {
	rec := httptest.NewRecorder()
	auth := handler.AuthChecker(realm, htpasswdTestFile)
	h := handlerFixture(t)
	chain := auth(h)
	request := httptest.NewRequest("GET", "/path", nil)
	request.SetBasicAuth("notvalid", "auth")
	chain.ServeHTTP(rec, request)
	assert.Equal(t, rec.Result().StatusCode, http.StatusUnauthorized)
	assertBodyContent(t, "401 Unauthorized\n", rec.Result().Body)
}

func TestAuthChecker_NotValidAuthFormat(t *testing.T) {
	rec := httptest.NewRecorder()
	auth := handler.AuthChecker(realm, htpasswdTestFile)
	h := handlerFixture(t)
	chain := auth(h)
	request := httptest.NewRequest("GET", "/path", nil)
	request.Header.Add("Authorization", "NOT_BASE64")
	chain.ServeHTTP(rec, request)
	assert.Equal(t, rec.Result().StatusCode, http.StatusUnauthorized)
	assertBodyContent(t, "401 Unauthorized\n", rec.Result().Body)
}