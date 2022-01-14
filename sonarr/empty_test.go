package sonarr_test

import (
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandler_Empty(t *testing.T) {
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")

	w := newResponseWriter()
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	handler.Empty(w, req)
	require.Equal(t, http.StatusForbidden, w.StatusCode)

	w = newResponseWriter()
	req.Header.Set("X-Api-Key", handler.APIKey)
	handler.Empty(w, req)
	require.Equal(t, http.StatusOK, w.StatusCode)
	assert.Equal(t, "[]", w.Response)
}
