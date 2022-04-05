package sonarr

import (
	"github.com/clambin/imdb-watchlist/watchlist"
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
