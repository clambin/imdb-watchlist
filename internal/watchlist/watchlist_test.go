package watchlist_test

import (
	"errors"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/clambin/imdb-watchlist/internal/watchlist/mocks"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
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
			s := watchlist.New(slog.Default(), r, metrics.NewRequestSummaryMetrics("", "", nil), tt.listIDs...)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v3/series", nil)
			s.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			if w.Code == http.StatusOK {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				assert.Equal(t, tt.body, w.Body.String())
			}
		})
	}
}

func TestServer_Movie(t *testing.T) {
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
				{IMDBId: "1", Type: imdb.Movie, Title: "some movie"},
				{IMDBId: "2", Type: imdb.Movie, Title: "some other movie"},
			}},
			wantStatusCode: http.StatusOK,
			body: `[{"title":"some movie","imdbId":"1"},{"title":"some other movie","imdbId":"2"}]
`,
		},
		{
			name:    "multiple",
			listIDs: []string{"1", "2"},
			responses: map[string]imdb.Watchlist{
				"1": {
					{IMDBId: "1", Type: imdb.Movie, Title: "foo"},
					{IMDBId: "2", Type: imdb.Movie, Title: "bar"},
				},
				"2": {
					{IMDBId: "1", Type: imdb.Movie, Title: "foo"},
					{IMDBId: "3", Type: imdb.Movie, Title: "snafu"},
				},
			},
			wantStatusCode: http.StatusOK,
			body: `[{"title":"foo","imdbId":"1"},{"title":"bar","imdbId":"2"},{"title":"snafu","imdbId":"3"}]
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

			s := watchlist.New(slog.Default(), r, metrics.NewRequestSummaryMetrics("", "", nil), tt.listIDs...)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v3/movie", nil)
			s.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			if w.Code == http.StatusOK {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				assert.Equal(t, tt.body, w.Body.String())
			}
		})
	}
}

func TestServer_Handle(t *testing.T) {
	s := watchlist.New(slog.Default(), nil, metrics.NewRequestSummaryMetrics("", "", nil))

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
