package watchlist

import "net/http"

const APIKeyHeader = "X-Api-Key"

func Authenticate(apiKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var authenticated bool
			for _, key := range req.Header.Values(APIKeyHeader) {
				if key == apiKey {
					authenticated = true
					break
				}
			}
			if !authenticated {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
