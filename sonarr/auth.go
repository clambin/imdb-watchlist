package sonarr

import (
	"net/http"
)

// AuthMiddleware checks that the request contains a valid API key
func (handler *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !authenticated(req, handler.APIKey) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("missing/invalid API key"))
			return
		}
		next.ServeHTTP(w, req)
	})
}

func authenticated(req *http.Request, apiKey string) bool {
	for _, key := range req.Header["X-Api-Key"] {
		if key == apiKey {
			return true
		}
	}
	return false
}
