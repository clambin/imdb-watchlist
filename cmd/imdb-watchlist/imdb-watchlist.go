package main

import (
	"context"
	"errors"
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
	"strings"
	"syscall"
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
		logger.Info("no API Key provided. generating a new one", "apikey", *apiKey)
	}

	reader := imdb.WatchlistFetcher{
		HTTPClient: &http.Client{
			Transport: httpclient.NewRoundTripper(httpclient.WithCache(httpclient.DefaultCacheTable, 15*time.Minute, time.Minute)),
		},
	}

	handler := watchlist.New(logger, reader, strings.Split(*listID, ",")...)
	prometheus.MustRegister(handler)

	tm := taskmanager.New(
		httpserver.New(
			*addr,
			middleware.RequestLogger(logger, slog.LevelInfo, middleware.RequestLogFormatterFunc(watchlist.FormatRequest))(
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

	return tm.Run(ctx)
}
