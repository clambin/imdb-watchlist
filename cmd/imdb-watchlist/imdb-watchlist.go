package main

import (
	"errors"
	"fmt"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/clambin/imdb-watchlist/version"
	"github.com/clambin/imdb-watchlist/watchlist"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
	"golang.org/x/exp/slog"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var (
		debug          bool
		port           int
		prometheusPort int
		listID         string
		apiKey         string
	)

	a := kingpin.New(filepath.Base(os.Args[0]), "imdb-watchlist")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&port)
	a.Flag("prometheus", "Prometheus listener port").Default("9090").IntVar(&prometheusPort)
	a.Flag("list", "IMDB List ID").Required().StringVar(&listID)
	a.Flag("apikey", "API Key").StringVar(&apiKey)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	var opts slog.HandlerOptions
	if debug {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}
	slog.SetDefault(slog.New(opts.NewTextHandler(os.Stdout)))

	slog.Info("imdb-watchlist starting", "version", version.BuildVersion)

	if apiKey == "" {
		apiKey = watchlist.GenerateKey()
		slog.Info("no API Key provided. generating a new one", "apikey", apiKey)
	}

	handler := watchlist.New(apiKey, &imdb.Client{
		HTTPClient: &http.Client{Transport: httpclient.NewRoundTripper(
			httpclient.WithCache{
				DefaultExpiry:   15 * time.Minute,
				CleanupInterval: time.Hour,
			},
		)},
		ListID: listID,
	})

	server := &http.Server{Addr: ":8080", Handler: handler.MakeRouter()}
	prometheus.MustRegister(handler)

	go func() {
		if err2 := server.ListenAndServe(); !errors.Is(err2, http.ErrServerClosed) {
			slog.Error("failed to start server", err2)
			panic(err2)
		}
	}()

	go runPrometheusServer(prometheusPort)

	<-shutdown.Chan()

	slog.Info("imdb-watchlist stopped")
}

func runPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus listener", err)
	}
}
