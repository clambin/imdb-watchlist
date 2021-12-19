package sonarr

import (
	"net/http"
)

func (handler *Handler) Empty(w http.ResponseWriter, req *http.Request) {
	if handler.handleAuth(w, req) == false {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`[]`))
}
