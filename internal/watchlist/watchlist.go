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
	http.Handler
	ListIDs []string
	reader  Reader
	metrics *middleware.PrometheusMetrics
	logger  *slog.Logger
}

var _ Reader = &imdb.WatchlistFetcher{}

type Reader interface {
	GetWatchlist(listID string) (entries imdb.Watchlist, err error)
}

func New(logger *slog.Logger, reader Reader, listIDs ...string) *Server {
	s := Server{
		ListIDs: listIDs,
		reader:  reader,
		metrics: middleware.NewPrometheusMetrics(middleware.PrometheusMetricsOptions{
			Application: "imdb-watchlist",
		}),
		logger: logger,
	}

	m := http.NewServeMux()
	m.Handle("GET /api/v3/series", s.metrics.Handle(http.HandlerFunc(s.Series)))
	m.Handle("GET /api/v3/movie", s.metrics.Handle(http.HandlerFunc(s.Movies)))
	m.HandleFunc("/api/v3/importList/action/getDevices", s.Empty)
	m.HandleFunc("/api/v3/qualityprofile", s.Empty)
	s.Handler = m

	return &s
}

func (s *Server) Series(w http.ResponseWriter, r *http.Request) {
	s.handleListRequest(w, r, "show", imdb.TVSeries, imdb.TVMiniSeries)
}

func (s *Server) Movies(w http.ResponseWriter, r *http.Request) {
	s.handleListRequest(w, r, "movie", imdb.Movie)
}

func (s *Server) handleListRequest(w http.ResponseWriter, _ *http.Request, mediaType string, entryType ...imdb.EntryType) {
	all, err := s.queryWatchLists(entryType...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(s.buildResponse(all, mediaType)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) queryWatchLists(entryTypes ...imdb.EntryType) (imdb.Watchlist, error) {
	var all imdb.Watchlist
	for _, listID := range s.ListIDs {
		watchlist, err := s.reader.GetWatchlist(listID)
		if err != nil {
			return nil, err
		}
		watchlist = watchlist.Filter(entryTypes...)
		s.logger.Debug("queried list", "listID", listID, "found", len(watchlist))
		all = append(all, watchlist...)
	}
	return Unique(all, func(v imdb.Entry) string { return v.IMDBId }), nil
}

// Entry represents an entry in the IMDB watchlist
type Entry struct {
	Title  string `json:"title"`
	IMDBId string `json:"imdbId"`
}

func (s *Server) buildResponse(imdbEntries []imdb.Entry, mediaType string) []Entry {
	entries := make([]Entry, len(imdbEntries))

	l := s.logger.With("type", mediaType)
	for i := range imdbEntries {
		l.Debug("imdb watchlist entry found", "title", imdbEntries[i].Title, "imdbId", imdbEntries[i].IMDBId)
		entries[i] = Entry{
			Title:  imdbEntries[i].Title,
			IMDBId: imdbEntries[i].IMDBId,
		}
	}

	return entries
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
