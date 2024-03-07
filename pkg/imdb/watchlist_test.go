package imdb_test

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWatchlist_Filter(t *testing.T) {
	baseWatchlist := imdb.Watchlist{
		{IMDBId: "1", Type: imdb.Movie, Title: "movie"},
		{IMDBId: "2", Type: imdb.TVSeries, Title: "series"},
		{IMDBId: "3", Type: imdb.TVMiniSeries, Title: "mini-series"},
	}

	tests := []struct {
		name       string
		entryTypes []imdb.EntryType
		want       imdb.Watchlist
	}{
		{
			name:       "single",
			entryTypes: []imdb.EntryType{imdb.Movie},
			want: imdb.Watchlist{
				{IMDBId: "1", Type: imdb.Movie, Title: "movie"},
			},
		},
		{
			name:       "multiple",
			entryTypes: []imdb.EntryType{imdb.TVSeries, imdb.TVMiniSeries},
			want: imdb.Watchlist{
				{IMDBId: "2", Type: imdb.TVSeries, Title: "series"},
				{IMDBId: "3", Type: imdb.TVMiniSeries, Title: "mini-series"},
			},
		},
		{
			name:       "missing",
			entryTypes: []imdb.EntryType{},
			want:       imdb.Watchlist{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, baseWatchlist.Filter(tt.entryTypes...))
		})
	}
}
