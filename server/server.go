package server

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/server"
	"github.com/clambin/imdb-watchlist/sonarr"
	"net/http"
	"time"
)

// Server runs an HTTP Server that provides the Sonarr endpoints, as well as a metrics endpoint for Prometheus
type Server struct {
	server.Server
}

// New creates a new Server
func New(port int, handler *sonarr.Handler) *Server {
	return &Server{
		Server: *server.NewWithHandlers(port, []server.Handler{
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
		}),
	}
}

// Run starts the HTTP server that provides the Sonarr endpoints
func (s *Server) Run(ctx context.Context) (err error) {
	ch := make(chan error)
	go func() {
		ch <- s.Server.Run()
	}()

	<-ctx.Done()

	_ = s.Server.Shutdown(5 * time.Second)
	if err = <-ch; errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}
