package server

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"net/http"
	"time"
)

// Run starts the HTTP server that provides the Sonarr endpoints
func Run(ctx context.Context, port int, handler *sonarr.Handler) (err error) {
	s := server.NewWithHandlers(port, []server.Handler{
		{
			Path:    "/api/v3/series",
			Handler: handler.AuthMiddleware(http.HandlerFunc(handler.Series)),
		},
		{
			Path:    "/api/v3/importList/action/getDevices",
			Handler: handler.AuthMiddleware(http.HandlerFunc(handler.Empty)),
		},
		{
			Path:    "api/v3/qualityprofile",
			Handler: handler.AuthMiddleware(http.HandlerFunc(handler.Empty)),
		},
	})

	ch := make(chan error)
	go func() {
		ch <- s.Run()
	}()

	<-ctx.Done()

	_ = s.Shutdown(5 * time.Second)
	if err = <-ch; errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}
