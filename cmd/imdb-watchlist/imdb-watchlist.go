package main

import (
	"context"
	"flag"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/go-common/taskmanager/httpserver"
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))

	slog.Info("imdb-watchlist starting", "version", version.BuildVersion)

	if *listID == "" {
		slog.Error("no IMDB List ID provided. Aborting.")
		os.Exit(1)
	}

	if *apiKey == "" {
		*apiKey = watchlist.GenerateKey()
		slog.Info("no API Key provided. generating a new one", "apikey", *apiKey)
	}

	handler := watchlist.New(*apiKey, &imdb.Fetcher{
		HTTPClient: &http.Client{
			Transport: httpclient.NewRoundTripper(httpclient.WithCache(httpclient.CacheTable{}, 15*time.Minute, time.Hour)),
		},
		ListID: *listID,
	})
	prometheus.MustRegister(handler)

	prom := http.NewServeMux()
	prom.Handle("/metrics", promhttp.Handler())

	tm := taskmanager.New(
		httpserver.New(*addr, handler.MakeRouter()),
		httpserver.New(*prometheusAddr, prom),
	)

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()
	if err := tm.Run(ctx); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}

	slog.Info("imdb-watchlist stopped")
}
