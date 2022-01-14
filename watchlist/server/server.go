package server

import (
	"net/http"
)

// Handler emulates an IMDB watchlist
type Handler struct {
	Fail     bool   // Fail any incoming call
	Response string // Response to return. If none is provided, defaults to ReferenceOutput
}

// Handle the incoming request
func (handler *Handler) Handle(w http.ResponseWriter, _ *http.Request) {
	if handler.Fail {
		http.Error(w, "server failure", http.StatusNotFound)
		return
	}
	/*
		if handler.Response == "" {
			handler.Response = ReferenceOutput
		}
	*/
	_, _ = w.Write([]byte(handler.Response))
}

// ReferenceOutput is a valid response
const ReferenceOutput = `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
4,tt4,,,,A TV miniseries,,tvMiniSeries,,,,,,,
`

// InvalidOutput is a syntactically invalid response
const InvalidOutput = `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
`

// HeaderMissing misses a mandatory column ("Const")
const HeaderMissing = `Position,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,,,,A Movie,,movie,,,,,,,
2,,,,A TV Series,,tvSeries,,,,,,,,
3,,,,A TV Special,,tvSpecial,,,,,,,
`
