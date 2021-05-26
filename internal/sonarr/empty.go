package sonarr

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (handler *Handler) Empty(w http.ResponseWriter, req *http.Request) {
	if handler.authenticate(req) == false {
		log.Warning("missing/invalid API key")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("missing/invalid API key"))
		return
	}

	log.WithField("url", req.URL.String()).Info("stubbed endpoint called")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`[]`))
}
