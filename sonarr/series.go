package sonarr

import (
	"encoding/json"
	"github.com/clambin/imdb-watchlist/watchlist"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (handler *Handler) Series(w http.ResponseWriter, req *http.Request) {

	if handler.handleAuth(w, req) == false {
		return
	}

	entries, err := handler.Client.Watchlist(handler.ListID, "tvSeries", "tvMiniSeries")

	if err != nil {
		log.WithError(err).Warning("failed to get IMDB Watchlist")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response []byte
	response, err = handler.buildResponse(entries)

	if err != nil {
		log.WithError(err).Error("failed to create response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(response)

}

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
		}).Info("found an entry")
	}

	return json.Marshal(sonarrEntries)
}
