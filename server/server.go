package server

import (
	"context"
	"github.com/clambin/gotools/metrics"
	"github.com/clambin/imdb-watchlist/sonarr"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Run(ctx context.Context, port int, handler *sonarr.Handler) {
	server := metrics.NewServerWithHandlers(port, []metrics.Handler{
		{
			Path:    "/api/v3/series",
			Handler: http.HandlerFunc(handler.Series),
		},
		{
			Path:    "/api/v3/importList/action/getDevices",
			Handler: http.HandlerFunc(handler.Empty),
		},
		{
			Path:    "api/v3/qualityprofile",
			Handler: http.HandlerFunc(handler.Empty),
		},
	})

	go func() {
		err := server.Run()
		if err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to start HTTP server")
		}
	}()

	<-ctx.Done()
	_ = server.Shutdown(5 * time.Second)
}
