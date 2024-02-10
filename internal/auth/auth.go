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
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
