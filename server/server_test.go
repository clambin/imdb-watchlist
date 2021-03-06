package server_test

import (
	"context"
	"fmt"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	"github.com/clambin/imdb-watchlist/pkg/watchlist/mocks"
	"github.com/clambin/imdb-watchlist/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	wl := &mocks.Reader{}
	handler := sonarr.New("12345", "ls1234")
	handler.Client = wl

	wl.
		On("GetByTypes", "tvSeries", "tvMiniSeries").
		Return([]watchlist.Entry{
			{IMDBId: "tt2", Title: "A TV Series"},
			{IMDBId: "tt4", Title: "A TV miniseries"},
		}, nil).Once()

	ctx, cancel := context.WithCancel(context.Background())

	s := server.New(0, handler)

	go func() {
		_ = s.Run(ctx)
	}()

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", s.Server.Port)

	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL + "/metrics")
		return err == nil && resp.StatusCode == http.StatusOK
	}, 500*time.Millisecond, 10*time.Millisecond)

	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v3/series", nil)
	req.Header["X-Api-Key"] = []string{"12345"}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `[{"title":"A TV Series","imdbId":"tt2"},{"title":"A TV miniseries","imdbId":"tt4"}]`, string(body))
	_ = resp.Body.Close()

	cancel()

	require.Never(t, func() bool {
		resp, err = http.Get(baseURL + "/metrics")
		return err == nil && resp.StatusCode == http.StatusOK
	}, 100*time.Millisecond, 10*time.Millisecond)
}
