package sonarr_test

import (
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/clambin/imdb-watchlist/watchlist/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Series(t *testing.T) {
	testServer := &mock.Handler{}
	server := httptest.NewServer(http.HandlerFunc(testServer.Handle))
	defer server.Close()

	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.Client.URL = server.URL

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Series(w, req)

	require.Equal(t, http.StatusOK, w.StatusCode)
	contentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, `[{"title":"A TV Series","imdbId":"tt2"},{"title":"A TV miniseries","imdbId":"tt4"}]`, w.Response)
}

func TestHandler_Series_NoAPIKey(t *testing.T) {
	testServer := &mock.Handler{}
	server := httptest.NewServer(http.HandlerFunc(testServer.Handle))
	defer server.Close()

	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.Client.URL = server.URL

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)

	handler.Series(w, req)
	assert.Equal(t, http.StatusForbidden, w.StatusCode)
}

func TestHandler_Series_FailedAPICall(t *testing.T) {
	testServer := &mock.Handler{}
	testServer.Fail = true
	server := httptest.NewServer(http.HandlerFunc(testServer.Handle))
	defer server.Close()

	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.Client.URL = server.URL

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Series(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
}
func TestHandler_Series_BadResponse(t *testing.T) {
	testServer := &mock.Handler{}
	testServer.Invalid = true
	server := httptest.NewServer(http.HandlerFunc(testServer.Handle))
	defer server.Close()

	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.Client.URL = server.URL

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Series(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
}

func TestHandler_Empty(t *testing.T) {
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)
	// req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Empty(w, req)
	assert.Equal(t, http.StatusForbidden, w.StatusCode)

	w = newResponseWriter()
	req, err = http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Empty(w, req)
	assert.Equal(t, http.StatusOK, w.StatusCode)
}

type ResponseWriter struct {
	StatusCode int
	Response   string
	Headers    http.Header
}

func newResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		StatusCode: 0,
		Response:   "",
		Headers:    make(map[string][]string),
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.Headers
}

func (w *ResponseWriter) Write(content []byte) (int, error) {
	w.Response += string(content)
	return len(content), nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}
