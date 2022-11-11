package server

import (
	"context"
	"fmt"
	"github.com/clambin/httpserver"
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	server *httpserver.Server
}

// New creates a new Server
func New(port int, handler *sonarr.Handler, r prometheus.Registerer) (s *Server, err error) {
	s = new(Server)
	options := []httpserver.Option{
		httpserver.WithPort{Port: port},
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
	}
	if r != nil {
		m := httpserver.NewMetrics("imdb-watchlist")
		m.Register(r)
		options = append(options, httpserver.WithMetrics{Metrics: m})
	}

	s.server, err = httpserver.New(options...)
	if err != nil {
		return nil, fmt.Errorf("handler: %w", err)
	}
	return s, nil
}

func (s *Server) GetPort() int {
	return s.server.GetPort()
}

func (s *Server) RunWithContext(ctx context.Context) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = s.server.Run()
		wg.Done()
	}()

	<-ctx.Done()
	_ = s.server.Shutdown(time.Minute)

	wg.Wait()
	return
}
