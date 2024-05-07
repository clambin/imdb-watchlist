package server

import (
	"context"
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/clambin/imdb-watchlist/internal/configuration"
	"github.com/clambin/imdb-watchlist/internal/testutils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	testServer := httptest.NewServer(testutils.IMDBServer("ls001"))
	t.Cleanup(testServer.Close)

	config := configuration.Configuration{
		Addr:     ":8080",
		PromAddr: ":9090",
		ListIDs:  []string{"ls001"},
		APIKey:   "1234567890",
		ImDbURL:  testServer.URL,
	}

	go func() {
		require.NoError(t, Run(ctx, config, nil, os.Stderr, "dev"))
	}()

	tests := []struct {
		name     string
		apiKey   string
		path     string
		wantCode int
		wantBody string
	}{
		/*
					{
						name:     "series",
						apiKey:   config.APIKey,
						path:     "/api/v3/series",
						wantCode: http.StatusOK,
						wantBody: `[{"title":"A TV Series","imdbId":"tt2"},{"title":"A TV Special","imdbId":"tt3"},{"title":"A TV miniseries","imdbId":"tt4"}]
			`,
					},
		*/
		{
			name:     "radarr list",
			apiKey:   config.APIKey,
			path:     "/api/v3/movie",
			wantCode: http.StatusOK,
			wantBody: `[{"title":"A Movie","imdbId":"tt1"}]
`,
		},
		{
			name:     "quality profile",
			apiKey:   config.APIKey,
			path:     "/api/v3/qualityprofile",
			wantCode: http.StatusOK,
			wantBody: `[]
`,
		},
		{
			name:     "devices",
			apiKey:   config.APIKey,
			path:     "/api/v3/importList/action/getDevices",
			wantCode: http.StatusOK,
			wantBody: `[]
`,
		},
		{
			name:     "missing api key",
			path:     "/api/v3/movie",
			wantCode: http.StatusUnauthorized,
			wantBody: `{ "error": "missing / wrong api key in header X-Api-Key" }
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, _ := http.NewRequest(http.MethodGet, "http://localhost"+config.Addr+tt.path, nil)
			if tt.apiKey != "" {
				req.Header.Set(auth.APIKeyHeader, tt.apiKey)
			}
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantCode, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestRun_Metrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	testServer := httptest.NewServer(testutils.IMDBServer("ls001"))
	t.Cleanup(testServer.Close)

	config := configuration.Configuration{
		Addr:     ":8080",
		PromAddr: ":9090",
		ListIDs:  []string{"ls001"},
		ImDbURL:  testServer.URL,
	}

	r := prometheus.NewRegistry()
	go func() {
		require.NoError(t, Run(ctx, config, r, os.Stderr, "dev"))
	}()

	_, err := http.Get("http://localhost" + config.Addr + "/api/v3/movie")
	assert.NoError(t, err)

	assert.NoError(t, testutil.CollectAndCompare(r, strings.NewReader(`
# HELP watchlist_http_requests_total total number of http requests
# TYPE watchlist_http_requests_total counter
watchlist_http_requests_total{code="200",method="GET",path="/api/v3/movie"} 1
`), "watchlist_http_requests_total"))
}
