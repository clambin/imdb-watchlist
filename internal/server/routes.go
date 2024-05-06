package server

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"log/slog"
	"net/http"
)

var stubbedRoutes = []string{
	"/api/v3/importList/action/getDevices",
	"/api/v3/qualityprofile",
	"/api/v3/languageprofile",
	"/api/v3/rootfolder",
	"/api/v3/tag",
}

func addRoutes(
	mux *http.ServeMux,
	fetcher WatchlistReader,
	listIDs []string,
	logger *slog.Logger,
) {
	mux.Handle("GET /api/v3/series", WatchlistHandler(
		fetcher,
		listIDs,
		[]imdb.EntryType{imdb.TVSeries, imdb.TVMiniSeries, imdb.TVSpecial},
		logger.With("handler", "imdbHandler", "type", "series"),
	))
	mux.Handle("GET /api/v3/movie", WatchlistHandler(
		fetcher,
		listIDs,
		[]imdb.EntryType{imdb.Movie},
		logger.With("handler", "imdbHandler", "type", "movie"),
	))
	for _, stubbedRoute := range stubbedRoutes {
		mux.Handle(stubbedRoute, EmptyHandler())
	}
}
