package sonarr

import (
	"encoding/json"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	log "github.com/sirupsen/logrus"
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
		log.WithError(err).Warning("failed to get IMDB GetByTypes")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

		log.WithFields(log.Fields{
			"title":  entry.Title,
			"imdbId": entry.IMDBId,
			"count":  len(sonarrEntries),
		}).Info("imdb watchlist entry found")
	}

	return json.Marshal(sonarrEntries)
}
