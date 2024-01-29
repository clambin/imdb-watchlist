package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/clambin/imdb-watchlist/internal/watchlist/mocks"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Handle(t *testing.T) {
	reader := mocks.NewReader(t)
	reader.On("GetWatchlist", imdb.TVSeries, imdb.TVMiniSeries).Return([]imdb.Entry{}, nil)

	s := watchlist.New(reader, slog.Default())

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(s)

	tests := []struct {
		name       string
		path       string
		statusCode int
	}{
		{
			name:       "series",
			path:       "/api/v3/series",
			statusCode: http.StatusOK,
		},
		{
			name:       "devices",
			path:       "/api/v3/importList/action/getDevices",
			statusCode: http.StatusOK,
		},
		{
			name:       "qualityProfile",
			path:       "/api/v3/qualityprofile",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "https://localhost"+tt.path, nil)

			s.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}

	count, err := testutil.GatherAndCount(reg)
	require.NoError(t, err)
	assert.Equal(t, 6, count)
}
