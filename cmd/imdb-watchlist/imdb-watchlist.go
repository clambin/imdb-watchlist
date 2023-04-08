package main

import (
	"context"
	"errors"
	"flag"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/clambin/imdb-watchlist/version"
	"github.com/clambin/imdb-watchlist/watchlist"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	debug          = flag.Bool("debug", false, "Log debug messages")
	addr           = flag.String("addr", ":8080", "Server address")
	prometheusAddr = flag.String("prometheus", ":9090", "Prometheus metrics address")
	listID         = flag.String("list", "", "IMDB List ID (required)")
	apiKey         = flag.String("apikey", "", "APIKey")
)

func main() {
	flag.Parse()

	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(opts.NewTextHandler(os.Stderr)))

	slog.Info("imdb-watchlist starting", "version", version.BuildVersion)

	if *listID == "" {
		slog.Error("no IMDB List ID provided. Aborting.")
		os.Exit(1)
	}

	if *apiKey == "" {
		*apiKey = watchlist.GenerateKey()
		slog.Info("no API Key provided. generating a new one", "apikey", *apiKey)
	}

	handler := watchlist.New(*apiKey, &imdb.Client{
		HTTPClient: &http.Client{Transport: httpclient.NewRoundTripper(
			httpclient.WithCache{
				DefaultExpiry:   15 * time.Minute,
				CleanupInterval: time.Hour,
			},
		)},
		ListID: *listID,
	})
	prometheus.MustRegister(handler)

	go func() {
		server := &http.Server{Addr: *addr, Handler: handler.MakeRouter()}
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "err", err)
			panic(err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(*prometheusAddr, nil); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start Prometheus listener", "err", err)
		}
	}()

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()
	<-ctx.Done()

	slog.Info("imdb-watchlist stopped")
}
