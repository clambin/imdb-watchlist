package sonarr

import (
	"encoding/json"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	"golang.org/x/exp/slog"
	"net/http"
)

// Series queries the IMDB watchlist and returns the contained TV series as subscribe series to Sonarr
func (handler *Handler) Series(w http.ResponseWriter, _ *http.Request) {
	entries, err := handler.Reader.GetByTypes("tvSeries", "tvMiniSeries")

	var response []byte
	if err == nil {
		response, err = handler.buildResponse(entries)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
}

// Entry represents an entry in the IMDB watchlist
type Entry struct {
	Title  string `json:"title"`
	IMDBId string `json:"imdbId"`
}

func (handler *Handler) buildResponse(entries []watchlist.Entry) (response []byte, err error) {
	sonarrEntries := make([]Entry, 0)

	for _, entry := range entries {
		sonarrEntries = append(sonarrEntries, Entry{
			Title:  entry.Title,
			IMDBId: entry.IMDBId,
		})

		slog.Info("imdb watchlist entry found",
			"title", entry.Title,
			"imdbId", entry.IMDBId,
			"count", len(sonarrEntries),
		)
	}

	return json.Marshal(sonarrEntries)
}
