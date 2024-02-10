package watchlist

import (
	"encoding/json"
	"github.com/clambin/go-common/httpserver/middleware"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
)

var _ prometheus.Collector = &Server{}

type Server struct {
	reader Reader
	http.Handler
	metrics *middleware.PrometheusMetrics
	logger  *slog.Logger
}

var _ Reader = &imdb.WatchlistFetcher{}

type Reader interface {
	GetWatchlist(validTypes ...imdb.EntryType) (entries []imdb.Entry, err error)
}

func New(reader Reader, logger *slog.Logger) *Server {
	s := Server{
		reader: reader,
		metrics: middleware.NewPrometheusMetrics(middleware.PrometheusMetricsOptions{
			Application: "imdb-watchlist",
		}),
		logger: logger,
	}
	s.Handler = s.makeRouter()

	return &s
}

func (s *Server) makeRouter() http.Handler {
	m := http.NewServeMux()
	m.Handle("GET /api/v3/series", s.metrics.Handle(http.HandlerFunc(s.Series)))
	m.HandleFunc("/api/v3/importList/action/getDevices", s.Empty)
	m.HandleFunc("/api/v3/qualityprofile", s.Empty)
	return m
}

func (s *Server) Series(w http.ResponseWriter, r *http.Request) {
	entries, err := s.getSeries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) Empty(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`[]`))
}

func (s *Server) Describe(ch chan<- *prometheus.Desc) {
	s.metrics.Describe(ch)
}

func (s *Server) Collect(ch chan<- prometheus.Metric) {
	s.metrics.Collect(ch)
}
