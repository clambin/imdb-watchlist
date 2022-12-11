package sonarr

import (
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

// Handler emulates a Sonarr server. It offers the necessary endpoints for a real Sonarr server to query it
// as a Sonarr Program List in Import Lists. When receiving a request for subscribed series (/api/v3/series endpoint),
// it will  query an IMDB watchlist and present it as a set of subscribed series.
type Handler struct {
	Reader    watchlist.Reader // queries an IMDB watchlist
	APIKey    string           // API key to expect from the calling Sonarr
	transport *httpclient.RoundTripper
}

var _ prometheus.Collector = &Handler{}

// New creates a new Handler
func New(apiKey, listID string) *Handler {
	transport := httpclient.NewRoundTripper(
		httpclient.WithCache{
			DefaultExpiry:   15 * time.Minute,
			CleanupInterval: time.Hour,
		})
	return &Handler{
		Reader: &watchlist.Client{
			HTTPClient: &http.Client{Transport: transport},
			ListID:     listID,
		},
		APIKey:    apiKey,
		transport: transport,
	}
}

func (handler *Handler) Describe(descs chan<- *prometheus.Desc) {
	handler.transport.Describe(descs)
}

func (handler *Handler) Collect(metrics chan<- prometheus.Metric) {
	handler.transport.Collect(metrics)
}
