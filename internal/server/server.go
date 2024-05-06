package server

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"log/slog"
	"net/http"
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
