package sonarr

import (
	"github.com/clambin/imdb-watchlist/watchlist"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	Client *watchlist.Client
	APIKey string
	ListID string
}

func New(apiKey, listID string) *Handler {
	return &Handler{
		Client: &watchlist.Client{},
		APIKey: apiKey,
		ListID: listID,
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
