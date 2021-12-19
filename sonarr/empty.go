package sonarr

import (
	log "github.com/sirupsen/logrus"
	"html"
	"net/http"
)

func (handler *Handler) Empty(w http.ResponseWriter, req *http.Request) {

	if handler.handleAuth(w, req) == false {
		return
	}

	log.WithField("url", html.EscapeString(req.URL.String())).Info("stubbed endpoint called")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`[]`))
}
