package imdb_test

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWatchlistFetcher_GetWatchlist(t *testing.T) {
	tests := []struct {
		name    string
		fail    bool
		wantErr assert.ErrorAssertionFunc
		want    imdb.Watchlist
	}{
		{
			name:    "pass",
			wantErr: assert.NoError,
			want: imdb.Watchlist{
				{IMDBId: "tt1", Type: imdb.Movie, Title: "A Movie"},
				{IMDBId: "tt2", Type: imdb.TVSeries, Title: "A TV Series"},
				{IMDBId: "tt3", Type: imdb.TVSpecial, Title: "A TV Special"},
				{IMDBId: "tt4", Type: imdb.TVMiniSeries, Title: "A TV miniseries"},
			},
		},
		{
			name:    "fail",
			fail:    true,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.fail {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, _ = w.Write([]byte(referenceOutput))
			}))
			defer s.Close()

			c := imdb.WatchlistFetcher{HTTPClient: http.DefaultClient, URL: s.URL}

			entries, err := c.GetWatchlist("1")
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, entries)
		})
	}
}

const referenceOutput = `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
4,tt4,,,,A TV miniseries,,tvMiniSeries,,,,,,,
`
