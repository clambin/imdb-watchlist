package auth

import "net/http"

const APIKeyHeader = "X-Api-Key"

func Authenticate(apiKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			for _, key := range req.Header.Values(APIKeyHeader) {
				if key == apiKey {
					next.ServeHTTP(w, req)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{ "error": "missing / wrong api key in header `+APIKeyHeader+`" }`, http.StatusUnauthorized)
		})
	}
}
