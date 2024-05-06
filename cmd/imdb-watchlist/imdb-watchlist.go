package main

import (
	"context"
	"fmt"
	"github.com/clambin/imdb-watchlist/internal/configuration"
	"github.com/clambin/imdb-watchlist/internal/server"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"os/signal"
	"syscall"
)

var (
	version = "change-me"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := configuration.GetConfiguration()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "invalid configuration: %s\n", err.Error())
		os.Exit(1)

	}
	if err = server.Run(ctx, cfg, prometheus.DefaultRegisterer, os.Stderr, version); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to start: %s\n", err.Error())
		os.Exit(1)
	}
}
