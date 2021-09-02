package server_test

import (
	"context"
	"github.com/clambin/imdb-watchlist/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/clambin/imdb-watchlist/watchlist/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	imdb := mock.Handler{}
	testImdb := httptest.NewServer(http.HandlerFunc(imdb.Handle))
	defer testImdb.Close()

	handler := sonarr.New("12345", "ls1234")
	handler.Client.URL = testImdb.URL

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		server.Run(ctx, 8080, handler)
	}()

	assert.Eventually(t, func() bool {
		resp, err := http.Get("http://localhost:8080/metrics")
		return err == nil && resp.StatusCode == http.StatusOK
	}, 500*time.Millisecond, 10*time.Millisecond)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v3/series", nil)
	req.Header["X-Api-Key"] = []string{"12345"}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `[{"title":"A TV Series","imdbId":"tt2"}]`, string(body))
	_ = resp.Body.Close()

	cancel()

	assert.Never(t, func() bool {
		resp, err = http.Get("http://localhost:8080/metrics")
		return err == nil && resp.StatusCode == http.StatusOK
	}, 100*time.Millisecond, 10*time.Millisecond)
}
