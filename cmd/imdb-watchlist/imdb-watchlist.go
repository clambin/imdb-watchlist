package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/imdb-watchlist/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/clambin/imdb-watchlist/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
	"golang.org/x/exp/slog"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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
	a.Flag("list", "IMDB GetByTypes ID").Required().StringVar(&listID)
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

	go runPrometheusServer(prometheusPort)

	if apiKey == "" {
		apiKey = sonarr.GenerateKey()
		slog.Info("no API Key provided. generating a new one", "apikey", apiKey)
	}

	h := sonarr.New(apiKey, listID)
	s, err := server.New(port, h)
	if err != nil {
		slog.Error("failed to start server", err)
		panic(err)
	}
	prometheus.MustRegister(s.HTTPServer, h)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err = s.RunWithContext(ctx); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start HTTP server", err)
			panic(err)
		}
		wg.Done()
	}()

	<-shutdown.Chan()
	cancel()
	wg.Wait()

	slog.Info("imdb-watchlist stopped")
}

func runPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus listener", err)
	}
}
