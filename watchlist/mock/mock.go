package mock

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	Fail    bool
	Invalid bool
}

func (handler *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	log.WithField("path", req.URL.Path).Debug("Handler")

	if handler.Fail {
		http.Error(w, "server failure", http.StatusNotFound)
		return
	}

	response := ReferenceOutput
	if handler.Invalid {
		response = InvalidOutput
	}
	_, _ = w.Write([]byte(response))
}

const ReferenceOutput = `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
`

const InvalidOutput = `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
`
