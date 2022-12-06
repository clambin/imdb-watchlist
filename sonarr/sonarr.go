package sonarr

import (
	"github.com/clambin/go-common/cache"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	"time"
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
		Client: &watchlist.Client{
			Caller: &httpclient.Cacher{
				Caller: &httpclient.BaseClient{},
				//Table:  client.CacheTable{},
				Cache: cache.New[string, []byte](15*time.Minute, time.Hour),
			},
			ListID: listID,
		},
		APIKey: apiKey,
	}
}
