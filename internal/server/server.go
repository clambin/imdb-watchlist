package server

import (
	"cmp"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"log/slog"
	"net/http"
	"slices"
)

type WatchlistReader interface {
	GetWatchlist(id string) (imdb.Watchlist, error)
}

func New(watchlistIDs []string, reader WatchlistReader, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, reader, watchlistIDs, logger)
	return mux
}

type Entry struct {
	Title  string `json:"title"`
	IMDBId string `json:"imdbId"`
}

func buildResponse(imdbEntries []imdb.Entry) []Entry {
	entries := make([]Entry, len(imdbEntries))

	for i := range imdbEntries {
		entries[i] = Entry{
			Title:  imdbEntries[i].Title,
			IMDBId: imdbEntries[i].IMDBId,
		}
	}

	return entries
}

func unique[V any, K cmp.Ordered](input []V, getKey func(V) K) []V {
	slices.SortFunc(input, func(a, b V) int {
		return cmp.Compare(getKey(a), getKey(b))
	})
	var last K
	entries := make([]V, 0, len(input))
	for _, e := range input {
		if key := getKey(e); key != last {
			entries = append(entries, e)
			last = key
		}
	}
	return entries
}
