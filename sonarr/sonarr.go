package sonarr

import (
	"net/http"
)

type Handler struct {
	HTTPClient *http.Client
	APIKey     string
	ListID     string
}

func New(apiKey, listID string) *Handler {
	return &Handler{HTTPClient: &http.Client{}, APIKey: apiKey, ListID: listID}
}

func (handler *Handler) authenticate(req *http.Request) bool {
	passedKeys := req.Header["X-Api-Key"]
	return len(passedKeys) > 0 && passedKeys[0] == handler.APIKey
}
