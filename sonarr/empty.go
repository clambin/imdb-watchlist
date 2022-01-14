package sonarr

import (
	"net/http"
)

// Empty handles endpoints that Handler does not support. This is mainly needed for api/v3/qualityprofile, which
// we don't support.
func (handler *Handler) Empty(w http.ResponseWriter, req *http.Request) {
	if handler.handleAuth(w, req) == false {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`[]`))
}
