//nolint:bodyclose
package www_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/www"
)

func TestApply(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("handler_write\n"))
	})
	var m1 www.Middleware = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("middleware1_write\n"))
			h.ServeHTTP(w, r)
		})
	}
	var m2 www.Middleware = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("middleware2_write\n"))
			h.ServeHTTP(w, r)
		})
	}
	a := www.Apply(h, m1, m2)
	rec := httptest.NewRecorder()
	a.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	expected := `middleware1_write
middleware2_write
handler_write
`
	body, err := ioutil.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, string(body))
}
