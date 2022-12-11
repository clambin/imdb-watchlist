package server

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/httpserver"
	"github.com/clambin/imdb-watchlist/sonarr"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	HTTPServer *httpserver.Server
}

// New creates a new Server
func New(port int, handler *sonarr.Handler) (s *Server, err error) {
	s = new(Server)
	s.HTTPServer, err = httpserver.New(
		httpserver.WithPort{Port: port},
		httpserver.WithMetrics{Application: "imdb-watchlist"},
		httpserver.WithHandlers{Handlers: []httpserver.Handler{
			{
				Path:    "/api/v3/series",
				Handler: handler.AuthMiddleware(http.HandlerFunc(handler.Series)),
			},
			{
				Path:    "/api/v3/importList/action/getDevices",
				Handler: handler.AuthMiddleware(http.HandlerFunc(handler.Empty)),
			},
			{
				Path:    "/api/v3/qualityprofile",
				Handler: handler.AuthMiddleware(http.HandlerFunc(handler.Empty)),
			},
		}},
	)
	if err != nil {
		return nil, fmt.Errorf("handler: %w", err)
	}
	return s, nil
}

func (s *Server) GetPort() int {
	return s.HTTPServer.GetPort()
}

func (s *Server) RunWithContext(ctx context.Context) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = s.HTTPServer.Serve()
		wg.Done()
	}()

	<-ctx.Done()
	_ = s.HTTPServer.Shutdown(time.Minute)

	wg.Wait()
	return
}
