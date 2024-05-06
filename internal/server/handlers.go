package server

import (
	"encoding/json"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/clambin/imdb-watchlist/pkg/unique"
	"log/slog"
	"net/http"
)

func WatchlistHandler(watcher WatchlistReader, listIDs []string, mediaTypes []imdb.EntryType, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var list imdb.Watchlist
		for _, listID := range listIDs {
			entries, err := watcher.GetWatchlist(listID)
			if err != nil {
				logger.Error("failed to get watchlist", "err", err, "listID", listID)
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
			list = append(list, entries.Filter(mediaTypes...)...)
		}
		list = unique.UniqueFunc(list, func(v imdb.Entry) string { return v.IMDBId })
		logger.Info("watchlist entries found", "list", list)
		writeResponse(w, http.StatusOK, buildResponse(list))
	})
}

func EmptyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeResponse[any](w, http.StatusOK, []any{})
	})
}

func writeResponse[T any](w http.ResponseWriter, statusCode int, response T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}
