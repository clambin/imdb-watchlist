package configuration

import (
	"errors"
	"flag"
	"strings"
)

var (
	debug    = flag.Bool("debug", false, "Log debug messages")
	addr     = flag.String("addr", ":8080", "Server address")
	promAddr = flag.String("prometheus", ":9090", "Prometheus metrics address")
	listID   = flag.String("list", "", "IMDB List ID(s) (required, comma-separated)")
	apiKey   = flag.String("apikey", "", "APIKey")
)

type Configuration struct {
	Debug    bool
	Addr     string
	PromAddr string
	ListIDs  []string
	APIKey   string
	ImDbURL  string
}

func GetConfiguration() (Configuration, error) {
	flag.Parse()
	if *listID == "" {
		return Configuration{}, errors.New("list is required")
	}
	return Configuration{
		Debug:    *debug,
		Addr:     *addr,
		PromAddr: *promAddr,
		ListIDs:  strings.Split(*listID, ","),
		APIKey:   *apiKey,
	}, nil
}
