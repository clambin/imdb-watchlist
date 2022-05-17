package sonarr

import (
	"net/http"
)

// AuthMiddleware checks that the request contains a valid API key
func (handler *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		passedKeys := req.Header["X-Api-Key"]
		authenticated := len(passedKeys) > 0 && passedKeys[0] == handler.APIKey

		if !authenticated {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("missing/invalid API key"))
			return
		}

		next.ServeHTTP(w, req)
	})
}
