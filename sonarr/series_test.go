package sonarr_test

import (
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/clambin/imdb-watchlist/watchlist/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandler_Series(t *testing.T) {
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.HTTPClient = httpstub.NewTestClient(mock.Serve)
	mock.ServerOutput = mock.ReferenceOutput

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Series(w, req)

	assert.Equal(t, http.StatusOK, w.StatusCode)
	contentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, `[{"title":"A TV Series","imdbId":"tt2"}]`, w.Response)
}

func TestHandler_Series_NoAPIKey(t *testing.T) {
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.HTTPClient = httpstub.NewTestClient(mock.Serve)
	mock.ServerOutput = mock.ReferenceOutput

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)

	handler.Series(w, req)
	assert.Equal(t, http.StatusForbidden, w.StatusCode)
}

func TestHandler_Series_FailedAPICall(t *testing.T) {
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.HTTPClient = httpstub.NewTestClient(Serve)

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Series(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
}

func TestHandler_Series_BadResponse(t *testing.T) {
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.HTTPClient = httpstub.NewTestClient(mock.Serve)
	mock.ServerOutput = ``

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Series(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
}

func Serve(_ *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
	}
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
