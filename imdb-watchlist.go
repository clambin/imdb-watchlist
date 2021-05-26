package main

import (
	"fmt"
	"github.com/clambin/imdb-watchlist/internal/sonarr"
	"github.com/clambin/imdb-watchlist/internal/version"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"path/filepath"
)

var (
	Debug  bool
	Port   int
	ListID string
	APIKey string
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "imdb-watchlist")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&Debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&Port)
	a.Flag("list", "IMDB Watchlist ID").Required().StringVar(&ListID)
	a.Flag("apikey", "API Key").StringVar(&APIKey)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if Debug {
		log.SetLevel(log.DebugLevel)
	}

	if APIKey == "" {
		APIKey = sonarr.GenerateKey()
		log.WithField("apikey", APIKey).Info("no API Key provided. generating a new one")
	}

	handler := sonarr.New(APIKey, ListID)
	r := &mux.Router{}
	r.Use(prometheusMiddleWare)
	r.Path("/metrics").Handler(promhttp.Handler())
	r.HandleFunc("/api/v3/series", handler.Series)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", 8080), r)
}

var (
	httpDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

func prometheusMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}
