package main

import (
	"errors"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/clambin/imdb-watchlist/version"
	"github.com/clambin/imdb-watchlist/watchlist"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	var (
		debug          bool
		addr           string
		prometheusAddr string
		listID         string
		apiKey         string
	)

	a := kingpin.New(filepath.Base(os.Args[0]), "imdb-watchlist")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&debug)
	a.Flag("addr", "API listener address").Default(":8080").StringVar(&addr)
	a.Flag("prometheus", "Prometheus listener address").Default(":9090").StringVar(&prometheusAddr)
	a.Flag("list", "IMDB List ID").Required().StringVar(&listID)
	a.Flag("apikey", "API Key").StringVar(&apiKey)

	if _, err := a.Parse(os.Args[1:]); err != nil {
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
	prometheus.MustRegister(handler)

	go func() {
		server := &http.Server{Addr: addr, Handler: handler.MakeRouter()}
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "err", err)
			panic(err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(prometheusAddr, nil); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start Prometheus listener", "err", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	slog.Info("imdb-watchlist stopped")
}
