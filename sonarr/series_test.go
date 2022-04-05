package sonarr_test

import (
	"errors"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/clambin/imdb-watchlist/watchlist"
	"github.com/clambin/imdb-watchlist/watchlist/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandler_Series(t *testing.T) {
	wl := &mocks.Reader{}
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.Client = wl

	wl.
		On("GetByTypes", "tvSeries", "tvMiniSeries").
		Return([]watchlist.Entry{
			{IMDBId: "tt2", Title: "A TV Series"},
			{IMDBId: "tt4", Title: "A TV miniseries"},
		}, nil).Once()

	w := newResponseWriter()
	req, _ := http.NewRequest(http.MethodGet, "", nil)

	handler.Series(w, req)
	require.Equal(t, http.StatusOK, w.StatusCode)

	contentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, `[{"title":"A TV Series","imdbId":"tt2"},{"title":"A TV miniseries","imdbId":"tt4"}]`, w.Response)

	mock.AssertExpectationsForObjects(t, wl)
}

func TestHandler_Series_FailedAPICall(t *testing.T) {
	wl := &mocks.Reader{}
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")
	handler.Client = wl

	wl.
		On("GetByTypes", "tvSeries", "tvMiniSeries").
		Return([]watchlist.Entry{}, errors.New("API call failed")).
		Once()

	w := newResponseWriter()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)

	handler.Series(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)

	mock.AssertExpectationsForObjects(t, wl)
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
