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

var (
	Debug  bool
	Port   int
	ListID string
	APIKey string
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "imdb-watchlist")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&Debug)
	a.Flag("port", "API listener port").Default("8080").IntVar(&Port)
	a.Flag("list", "IMDB Watchlist ID").Required().StringVar(&ListID)
	a.Flag("apikey", "API Key").StringVar(&APIKey)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("version", version.BuildVersion).Info("imdb-watchlist starting")

	if APIKey == "" {
		APIKey = sonarr.GenerateKey()
		log.WithField("apikey", APIKey).Info("no API Key provided. generating a new one")
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		server.Run(ctx, Port, sonarr.New(APIKey, ListID))
		wg.Done()
	}()

	<-shutdown.Chan()
	cancel()
	wg.Wait()

	log.Info("imdb-watchlist stopped")
}
