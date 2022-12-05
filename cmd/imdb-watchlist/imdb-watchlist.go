package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/httpserver"
	"github.com/clambin/imdb-watchlist/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/clambin/imdb-watchlist/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
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

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("version", version.BuildVersion).Info("imdb-watchlist starting")

	go runPrometheusServer(prometheusPort)

	if apiKey == "" {
		apiKey = sonarr.GenerateKey()
		log.WithField("apikey", apiKey).Info("no API Key provided. generating a new one")
	}

	metrics := httpserver.NewAvgMetrics("imdb-watchlist")
	prometheus.MustRegister(metrics)

	s, err := server.New(port, sonarr.New(apiKey, listID), metrics)
	if err != nil {
		log.WithError(err).Fatal("failed to start server")
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err = s.RunWithContext(ctx); err != nil {
			log.WithError(err).Fatal("failed to start HTTP server")
		}
		wg.Done()
	}()

	<-shutdown.Chan()
	cancel()
	wg.Wait()

	log.Info("imdb-watchlist stopped")
}

func runPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Fatal("failed to start Prometheus listener")
	}
}
