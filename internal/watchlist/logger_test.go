package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func TestFormatRequest(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/foo", nil)
	req.RemoteAddr = "localhost:12345"
	req.Header.Set("User-Agent", "bar")
	attrs := watchlist.FormatRequest(req, http.StatusOK, 10*time.Millisecond)

	v := slog.GroupValue(attrs...)
	assert.Equal(t, "[path=/foo method=GET source=localhost:12345 agent=bar code=200 latency=10ms]", v.String())
}
