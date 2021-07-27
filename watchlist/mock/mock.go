package mock

import (
	"bytes"
	"io"
	"net/http"
)

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

const EmptyOutput = ``

var ServerOutput string

func Serve(_ *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(ServerOutput)),
	}
}
