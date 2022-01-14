package sonarr

import (
	"github.com/clambin/imdb-watchlist/watchlist"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Handler emulates a Sonarr server. It offers the necessary endpoints for a real Sonarr server to query it
// as a Sonarr Program List in Import Lists. When receiving a request for subscribed series (/api/v3/series endpoint),
// it will  query an IMDB watchlist and present it as a set of subscribed series.
type Handler struct {
	Client watchlist.Reader // queries an IMDB watchlist
	APIKey string           // API key to expect from the calling Sonarr
}

// New creates a new Handler
func New(apiKey, listID string) *Handler {
	return &Handler{
		Client: &watchlist.Client{ListID: listID},
		APIKey: apiKey,
	}
}

func (handler *Handler) authenticate(req *http.Request) bool {
	passedKeys := req.Header["X-Api-Key"]
	return len(passedKeys) > 0 && passedKeys[0] == handler.APIKey
}

func (handler *Handler) handleAuth(w http.ResponseWriter, req *http.Request) bool {
	if handler.authenticate(req) == false {
		log.Warning("missing/invalid API key")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("missing/invalid API key"))
		return false
	}
	return true
}
