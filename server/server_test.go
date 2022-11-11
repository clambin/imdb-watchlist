package server_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	"github.com/clambin/imdb-watchlist/pkg/watchlist/mocks"
	"github.com/clambin/imdb-watchlist/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update golden images")

func TestRun(t *testing.T) {
	handler := sonarr.New("12345", "ls1234")
	wl := mocks.NewReader(t)
	handler.Client = wl

	wl.
		On("GetByTypes", "tvSeries", "tvMiniSeries").
		Return([]watchlist.Entry{
			{IMDBId: "tt2", Title: "A TV Series"},
			{IMDBId: "tt4", Title: "A TV miniseries"},
		}, nil).Once()

	ctx, cancel := context.WithCancel(context.Background())

	r := prometheus.NewRegistry()
	s, err := server.New(0, handler, r)
	require.NoError(t, err)

	go func() {
		_ = s.RunWithContext(ctx)
	}()

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", s.GetPort())

	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL + "/api/v3/series")
		if err == nil {
			_ = resp.Body.Close()
		}
		return err == nil && resp.StatusCode == http.StatusForbidden
	}, time.Second, time.Millisecond)

	endpoints := []string{
		"/api/v3/series",
		"/api/v3/importList/action/getDevices",
		"/api/v3/qualityprofile",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, baseURL+endpoint, nil)
			req.Header["X-Api-Key"] = []string{"12345"}
			var resp *http.Response
			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			gp := path.Join("testdata", strings.ToLower(t.Name()+".golden"))

			if *update {
				err = os.WriteFile(gp, body, 0644)
				require.NoError(t, err)
			}

			var golden []byte
			golden, err = os.ReadFile(gp)
			require.NoError(t, err)
			assert.Equal(t, string(golden), string(body))
		})
	}

	cancel()
}

func TestServer_BadPort(t *testing.T) {
	handler := sonarr.New("12345", "ls1234")
	_, err := server.New(-1, handler, nil)
	require.Error(t, err)
}
