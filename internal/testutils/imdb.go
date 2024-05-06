package testutils

import "net/http"

const watchlistBody = `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
4,tt4,,,,A TV miniseries,,tvMiniSeries,,,,,,,
`

func IMDBServer(listID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/list/" + listID + "/export":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(watchlistBody))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
