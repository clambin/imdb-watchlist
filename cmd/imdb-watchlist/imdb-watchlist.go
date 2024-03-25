package main

import (
	"errors"
	"flag"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/middleware"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	version = "change-me"

	debug          = flag.Bool("debug", false, "Log debug messages")
	addr           = flag.String("addr", ":8080", "Server address")
	prometheusAddr = flag.String("prometheus", ":9090", "Prometheus metrics address")
	listID         = flag.String("list", "", "IMDB List ID(s) (required, comma-separated)")
	apiKey         = flag.String("apikey", "", "APIKey")
)

func main() {
	flag.Parse()
	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))
	if err := Main(logger); err != nil {
		logger.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func Main(logger *slog.Logger) error {
	if *listID == "" {
		return errors.New("no IMDB List ID provided")
	}

	//TODO: not exactly secure. Create a separate tool to generate a key?
	if *apiKey == "" {
		*apiKey, _ = auth.GenerateKey()
		logger.Warn("no API Key provided. generating a new one", "apikey", *apiKey)
	}

	go func() {
		if err := http.ListenAndServe(*prometheusAddr, promhttp.Handler()); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start prometheus server", "err", err)
			panic(err)
		}
	}()

	clientMetrics := metrics.NewRequestSummaryMetrics("watchlist", "client", nil)
	prometheus.MustRegister(clientMetrics)

	reader := imdb.WatchlistFetcher{
		HTTPClient: &http.Client{
			Transport: roundtripper.New(
				roundtripper.WithCache(roundtripper.DefaultCacheTable, 15*time.Minute, time.Minute),
				roundtripper.WithRequestMetrics(clientMetrics),
			),
			Timeout: 10 * time.Second,
		},
	}

	serverMetrics := metrics.NewRequestSummaryMetrics("watchlist", "", nil)
	prometheus.MustRegister(serverMetrics)

	logger.Info("imdb-watchlist starting", "version", version)
	defer logger.Info("imdb-watchlist stopped")

	err := http.ListenAndServe(*addr,
		middleware.RequestLogger(logger, slog.LevelInfo, middleware.RequestLogFormatterFunc(watchlist.FormatRequest))(
			auth.Authenticate(*apiKey)(
				watchlist.New(logger, reader, serverMetrics, strings.Split(*listID, ",")...),
			),
		),
	)

	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start imdb-watchlist", "err", err)
		return err
	}
	return nil
}
