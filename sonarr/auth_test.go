package sonarr_test

import (
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	handler := sonarr.New(sonarr.GenerateKey(), "ls001")

	h := handler.AuthMiddleware(http.HandlerFunc(func(w2 http.ResponseWriter, _ *http.Request) {
		w2.WriteHeader(http.StatusOK)
	}))

	testCases := []struct {
		key        string
		statusCode int
	}{
		{key: handler.APIKey, statusCode: http.StatusOK},
		{key: "bad key", statusCode: http.StatusForbidden},
		{key: "", statusCode: http.StatusForbidden},
	}

	for _, testCase := range testCases {
		w := newResponseWriter()
		req, _ := http.NewRequest(http.MethodGet, "", nil)
		if testCase.key != "" {
			req.Header.Set("X-Api-Key", testCase.key)
		}

		h.ServeHTTP(w, req)

		assert.Equal(t, testCase.statusCode, w.StatusCode, testCase.key)
	}
}
