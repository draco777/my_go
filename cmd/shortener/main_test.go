package main

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", handleGet)
	r.Post("/", handlePost)
	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {

	data := []byte(body)
	r := bytes.NewReader(data)
	req, err := http.NewRequest(method, ts.URL+path, r)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {

	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/0", "")
	resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Плохой запрос\n", body)

	resp, body = testRequest(t, ts, "POST", "/", "my.ru")
	resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/0", body)

}
