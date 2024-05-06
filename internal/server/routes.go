package server

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	fetcher Reader,
	listIDs []string,
	logger *slog.Logger,
) {
	mux.Handle("GET /api/v3/series", IMDBHandler(
		fetcher,
		listIDs,
		[]imdb.EntryType{imdb.TVSeries, imdb.TVMiniSeries, imdb.TVSpecial},
		logger.With("handler", "imdbHandler", "type", "series"),
	))
	mux.Handle("GET /api/v3/movie", IMDBHandler(
		fetcher,
		listIDs,
		[]imdb.EntryType{imdb.Movie},
		logger.With("handler", "imdbHandler", "type", "series"),
	))
	mux.Handle("/api/v3/importList/action/getDevices", EmptyHandler())
	mux.Handle("/api/v3/qualityprofile", EmptyHandler())
}
