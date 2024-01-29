package watchlist_test

import (
	"errors"
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
	r := mocks.NewReader(t)
	s := watchlist.New(r, slog.Default())

	tests := []struct {
		name    string
		entries []imdb.Entry
		err     error
		pass    bool
		body    string
	}{
		{
			name:    "empty",
			entries: []imdb.Entry{},
			err:     nil,
			pass:    true,
			body:    "[]\n",
		},
		{
			name: "not empty",
			entries: []imdb.Entry{
				{IMDBId: "1", Type: "tvSeries", Title: "some series"},
				{IMDBId: "2", Type: "tvMiniSeries", Title: "some miniseries"},
			},
			err:  nil,
			pass: true,
			body: "[{\"title\":\"some series\",\"imdbId\":\"1\"},{\"title\":\"some miniseries\",\"imdbId\":\"2\"}]\n",
		},
		{
			name: "error",
			err:  errors.New("fail"),
			pass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.On("GetWatchlist", imdb.TVSeries, imdb.TVMiniSeries).Return(tt.entries, tt.err).Once()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v3/series", nil)

			s.Series(w, req)

			if !tt.pass {
				assert.NotEqual(t, http.StatusOK, w.Code)
				return
			}

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.body, w.Body.String())
		})
	}
}
