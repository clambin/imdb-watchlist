package watchlist

import (
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

var _ chiMiddleware.LogFormatter = &Logger{}

type Logger struct {
	logger *slog.Logger
}

func (l Logger) NewLogEntry(r *http.Request) chiMiddleware.LogEntry {
	return &LogEntry{logger: l.logger, request: r}
}

var _ chiMiddleware.LogEntry = &LogEntry{}

type LogEntry struct {
	logger  *slog.Logger
	request *http.Request
}

func (l LogEntry) Write(status, _ int, _ http.Header, elapsed time.Duration, _ interface{}) {
	l.logger.Info("request processed",
		slog.Group("request",
			slog.String("from", l.request.RemoteAddr),
			slog.String("path", l.request.URL.Path),
			slog.String("method", l.request.Method),
			slog.Int("status", status),
			//slog.Int("responseSize", bytes),
			slog.Duration("elapsed", elapsed),
		),
	)
}

func (l LogEntry) Panic(v interface{}, _ []byte) {
	chiMiddleware.PrintPrettyStack(v)
}
