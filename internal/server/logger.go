package server

import (
	"log/slog"
	"net/http"
	"time"
)

func FormatRequest(r *http.Request, statusCode int, latency time.Duration) []slog.Attr {
	return []slog.Attr{
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
		slog.String("source", r.RemoteAddr),
		slog.String("agent", r.UserAgent()),
		slog.Int("code", statusCode),
		slog.Duration("latency", latency),
	}
}
