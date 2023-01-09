package watchlist

import (
	"encoding/json"
	"github.com/clambin/go-common/httpserver/middleware"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-http-utils/headers"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"net/http"
)

type Server struct {
	APIKey  string
	Reader  Reader
	metrics *middleware.PrometheusMetrics
}

var _ prometheus.Collector = &Server{}

// Reader interface reads an IMDB watchlist
//
//go:generate mockery --name Reader
type Reader interface {
	ReadByTypes(validTypes ...string) (entries []imdb.Entry, err error)
}

var _ Reader = &imdb.Client{}

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

	r.Use(chiMiddleware.RequestLogger(&Logger{logger: slog.Default()}))
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

	w.Header().Set(headers.ContentType, "application/json")
	if err = json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) Empty(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(headers.ContentType, "application/json")
	_, _ = w.Write([]byte(`[]`))
}

func (s *Server) Describe(descs chan<- *prometheus.Desc) {
	s.metrics.Describe(descs)
}

func (s *Server) Collect(metrics chan<- prometheus.Metric) {
	s.metrics.Collect(metrics)
}