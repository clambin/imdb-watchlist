package watchlist_test

import (
	"bytes"
	"errors"
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/clambin/imdb-watchlist/internal/watchlist/mocks"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Series(t *testing.T) {
	tests := []struct {
		name           string
		listIDs        []string
		responses      map[string]imdb.Watchlist
		err            error
		wantStatusCode int
		body           string
	}{
		{
			name:           "empty",
			listIDs:        []string{"1"},
			responses:      map[string]imdb.Watchlist{"1": {}},
			wantStatusCode: http.StatusOK,
			body:           "[]\n",
		},
		{
			name:    "single",
			listIDs: []string{"1"},
			responses: map[string]imdb.Watchlist{"1": {
				{IMDBId: "1", Type: imdb.TVSeries, Title: "some series"},
				{IMDBId: "2", Type: imdb.TVMiniSeries, Title: "some miniseries"},
			}},
			wantStatusCode: http.StatusOK,
			body: `[{"title":"some series","imdbId":"1"},{"title":"some miniseries","imdbId":"2"}]
`,
		},
		{
			name:    "multiple",
			listIDs: []string{"1", "2"},
			responses: map[string]imdb.Watchlist{
				"1": {
					{IMDBId: "1", Type: imdb.TVSeries, Title: "some series"},
					{IMDBId: "2", Type: imdb.TVMiniSeries, Title: "some miniseries"},
				},
				"2": {
					{IMDBId: "1", Type: imdb.TVSeries, Title: "some series"},
					{IMDBId: "3", Type: imdb.TVMiniSeries, Title: "some other miniseries"},
				},
			},
			wantStatusCode: http.StatusOK,
			body: `[{"title":"some series","imdbId":"1"},{"title":"some miniseries","imdbId":"2"},{"title":"some other miniseries","imdbId":"3"}]
`,
		},
		{
			name:           "error",
			listIDs:        []string{"1"},
			responses:      map[string]imdb.Watchlist{"1": nil},
			err:            errors.New("fail"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewReader(t)
			for id, responses := range tt.responses {
				r.EXPECT().GetWatchlist(id).Return(responses, tt.err).Once()
			}
			s := watchlist.New(slog.Default(), r, tt.listIDs...)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v3/series", nil)
			s.Series(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			if w.Code == http.StatusOK {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				assert.Equal(t, tt.body, w.Body.String())
			}
		})
	}
}

func TestServer_Handle(t *testing.T) {
	s := watchlist.New(slog.Default(), nil)

	tests := []struct {
		name           string
		listIDs        []string
		responses      map[string]imdb.Watchlist
		path           string
		method         string
		wantStatusCode int
		want           string
	}{
		{
			name:           "series - wrong method",
			path:           "/api/v3/series",
			method:         http.MethodPost,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "devices",
			path:           "/api/v3/importList/action/getDevices",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "qualityProfile",
			path:           "/api/v3/qualityprofile",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, "https://localhost"+tt.path, nil)

			s.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
		})
	}
}

func TestServer_Collect(t *testing.T) {
	r := mocks.NewReader(t)
	r.EXPECT().GetWatchlist("1").Return(nil, nil).Once()
	s := watchlist.New(slog.Default(), r, "1")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v3/series", nil)
	s.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.NoError(t, testutil.CollectAndCompare(s, bytes.NewBufferString(`
# HELP http_requests_total Total number of http requests
# TYPE http_requests_total counter
http_requests_total{code="200",handler="imdb-watchlist",method="GET",path="/api/v3/series"} 1
`), "http_requests_total"))
}
