package handler_test

import (
	"github.com/eloylp/go-serve/handler"
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
