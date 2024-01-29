package watchlist

import (
	"encoding/json"
	"github.com/clambin/go-common/httpserver/middleware"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
)

type Server struct {
	APIKey  string
	Reader  Reader
	metrics *middleware.PrometheusMetrics
}

var _ prometheus.Collector = &Server{}

// Reader interface reads an IMDB watchlist
type Reader interface {
	ReadByTypes(validTypes ...imdb.EntryType) (entries []imdb.Entry, err error)
}

var _ Reader = &imdb.Fetcher{}

func New(apiKey string, reader Reader) *Server {
	s := Server{
		APIKey: apiKey,
		Reader: reader,
		metrics: middleware.NewPrometheusMetrics(middleware.PrometheusMetricsOptions{
			Application: "imdb-watchlist",
		}),
	}

	return &s
}

func (s *Server) MakeRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(slog.Default(), slog.LevelInfo, middleware.DefaultRequestLogFormatter))
	r.Use(Authenticate(s.APIKey))
	r.Use(s.metrics.Handle)

	r.Get("/api/v3/series", s.Series)
	r.Get("/api/v3/importList/action/getDevices", s.Empty)
	r.Get("/api/v3/qualityprofile", s.Empty)
	return r
}

func (s *Server) Series(w http.ResponseWriter, _ *http.Request) {
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
