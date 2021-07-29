package sonarr

import (
	"github.com/clambin/imdb-watchlist/watchlist"
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
