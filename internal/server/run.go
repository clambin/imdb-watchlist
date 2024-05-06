package server

import (
	"context"
	"errors"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/middleware"
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/clambin/imdb-watchlist/internal/configuration"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func Run(ctx context.Context, config configuration.Configuration, registerer prometheus.Registerer, logOutput io.Writer, version string) error {
	var opts slog.HandlerOptions
	if config.Debug {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(logOutput, &opts))
	logger.Info("imdb-watchlist starting", "version", version)

	promErr := make(chan error)
	go func() {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		promErr <- runHTTPServer(ctx, config.PromAddr, m)
	}()

	f := imdb.WatchlistFetcher{
		HTTPClient: http.DefaultClient,
		URL:        config.ImDbURL,
	}

	h := New(config.ListIDs, f, logger)
	if config.APIKey != "" {
		h = auth.Authenticate(config.APIKey)(h)
	}
	if registerer != nil {
		m := metrics.NewRequestSummaryMetrics("watchlist", "", nil)
		registerer.MustRegister(m)
		h = middleware.WithRequestMetrics(m)(h)
	}
	h = middleware.RequestLogger(logger, slog.LevelDebug, middleware.RequestLogFormatterFunc(FormatRequest))(h)

	serverErr := make(chan error)
	go func() {
		serverErr <- runHTTPServer(ctx, config.Addr, h)
	}()

	err := errors.Join(<-serverErr, <-promErr)
	logger.Info("imdb-watchlist stopped")
	return err
}

func runHTTPServer(ctx context.Context, addr string, handler http.Handler) error {
	httpServer := &http.Server{Addr: addr, Handler: handler}
	errCh := make(chan error)
	go func() {
		err := httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		errCh <- err
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel2()
	err := httpServer.Shutdown(ctx2)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}
