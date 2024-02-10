package auth_test

import (
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	const validKey = "1234"
	tests := []struct {
		name     string
		apiKey   string
		wantCode int
	}{
		{
			name:     "valid",
			apiKey:   validKey,
			wantCode: http.StatusOK,
		},
		{
			name:     "missing",
			wantCode: http.StatusForbidden,
		},
		{
			name:     "invalid",
			apiKey:   validKey + "5678",
			wantCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, _ := http.NewRequest(http.MethodGet, "", nil)
			if tt.apiKey != "" {
				r.Header.Set(auth.APIKeyHeader, tt.apiKey)
			}
			w := httptest.NewRecorder()

			auth.Authenticate(validKey)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
