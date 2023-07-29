package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/clambin/imdb-watchlist/internal/watchlist/mocks"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_MakeRouter(t *testing.T) {
	reader := mocks.NewReader(t)
	reader.On("ReadByTypes", imdb.TVSeries, imdb.TVMiniSeries).Return([]imdb.Entry{}, nil)

	s := watchlist.New("1234", reader)
	r := s.MakeRouter()

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(s)

	tests := []struct {
		name       string
		path       string
		apiKey     string
		statusCode int
	}{
		{
			name:       "series",
			path:       "/api/v3/series",
			apiKey:     "1234",
			statusCode: http.StatusOK,
		},
		{
			name:       "devices",
			path:       "/api/v3/importList/action/getDevices",
			apiKey:     "1234",
			statusCode: http.StatusOK,
		},
		{
			name:       "qualityProfile",
			path:       "/api/v3/qualityprofile",
			apiKey:     "1234",
			statusCode: http.StatusOK,
		},
		{
			name:       "missing key",
			path:       "/api/v3/series",
			statusCode: http.StatusForbidden,
		},
		{
			name:       "wrong key",
			path:       "/api/v3/series",
			apiKey:     "4321",
			statusCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "https://localhost"+tt.path, nil)
			if tt.apiKey != "" {
				req.Header.Set(watchlist.APIKeyHeader, tt.apiKey)
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}

	count, err := testutil.GatherAndCount(reg)
	require.NoError(t, err)
	assert.Equal(t, 6, count)
}
