package sonarr

import (
	"encoding/json"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (handler *Handler) Series(w http.ResponseWriter, req *http.Request) {

	if handler.authenticate(req) == false {
		log.Warning("missing/invalid API key")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("missing/invalid API key"))
		return
	}

	entries, err := watchlist.Get(handler.HTTPClient, handler.ListID, "tvSeries")

	if err != nil {
		log.WithError(err).Warning("failed to get IMDB Watchlist")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(handler.buildResponse(entries)))

}

type Entry struct {
	Title  string `json:"title"`
	IMDBId string `json:"imdbId"`
}

func (handler *Handler) buildResponse(entries []watchlist.Entry) (response string) {
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

	var output []byte
	var err error

	if output, err = json.Marshal(sonarrEntries); err != nil {
		log.WithError(err).Error("unable to build API response")
		output = []byte{}
	}

	return string(output)
}
