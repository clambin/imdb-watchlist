package main

import (
	"context"
	"github.com/clambin/imdb-watchlist/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/clambin/imdb-watchlist/version"
	log "github.com/sirupsen/logrus"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var (
		debug  bool
		port   int
		listID string
		apiKey string
	)

	a := kingpin.New(filepath.Base(os.Args[0]), "imdb-watchlist")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&port)
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

	if apiKey == "" {
		apiKey = sonarr.GenerateKey()
		log.WithField("apikey", apiKey).Info("no API Key provided. generating a new one")
	}

	s := server.New(port, sonarr.New(apiKey, listID))

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err = s.Run(ctx); err != nil {
			log.WithError(err).Fatal("failed to start HTTP server")
		}
		wg.Done()
	}()

	<-shutdown.Chan()
	cancel()
	wg.Wait()

	log.Info("imdb-watchlist stopped")
}
