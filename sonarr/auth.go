package sonarr

import (
	"net/http"
)

func (handler *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		passedKeys := req.Header["X-Api-Key"]
		authenticated := len(passedKeys) > 0 && passedKeys[0] == handler.APIKey

		if authenticated == false {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("missing/invalid API key"))
			return
		}

		next.ServeHTTP(w, req)
	})
}
