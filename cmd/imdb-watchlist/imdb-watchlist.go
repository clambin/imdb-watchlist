package main

import (
	"context"
	"flag"
	"github.com/clambin/go-common/httpclient"
	"github.com/clambin/go-common/httpserver/middleware"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/go-common/taskmanager/httpserver"
	promserver "github.com/clambin/go-common/taskmanager/prometheus"
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	version = "change-me"

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
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))

	if *listID == "" {
		logger.Error("no IMDB List ID provided. Aborting.")
		os.Exit(1)
	}

	//TODO: not exactly secure. Create a separate tool to generate a key?
	if *apiKey == "" {
		*apiKey, _ = auth.GenerateKey()
		logger.Info("no API Key provided. generating a new one", "apikey", *apiKey)
	}

	r := imdb.WatchlistFetcher{
		HTTPClient: &http.Client{
			Transport: httpclient.NewRoundTripper(httpclient.WithCache(httpclient.DefaultCacheTable, 15*time.Minute, time.Hour)),
		},
		ListID: *listID,
	}

	handler := watchlist.New(r, logger.With("component", "watchlist"))
	prometheus.MustRegister(handler)

	tm := taskmanager.New(
		httpserver.New(
			*addr,
			middleware.RequestLogger(logger, slog.LevelInfo, middleware.RequestLogFormatterFunc(formatRequest))(
				auth.Authenticate(*apiKey)(
					handler,
				),
			),
		),
		promserver.New(promserver.WithAddr(*prometheusAddr)),
	)

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	logger.Info("imdb-watchlist starting", "version", version)
	defer logger.Info("imdb-watchlist stopped")

	if err := tm.Run(ctx); err != nil {
		logger.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func formatRequest(r *http.Request, statusCode int, latency time.Duration) []slog.Attr {
	return []slog.Attr{
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
		slog.String("source", r.RemoteAddr),
		slog.String("agent", r.UserAgent()),
		slog.Int("code", statusCode),
		slog.Duration("latency", latency),
	}
}
