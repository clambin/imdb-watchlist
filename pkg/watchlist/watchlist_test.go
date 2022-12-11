package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	name       string
	fail       bool
	response   string
	validTypes []string
	output     []string
	pass       bool
}

var GetTestCases = []TestCase{
	{
		name:       "tvSeries",
		fail:       false,
		response:   ReferenceOutput,
		validTypes: []string{"tvSeries"},
		output:     []string{"tt2"},
		pass:       true,
	},
	{
		name:       "movie",
		fail:       false,
		response:   ReferenceOutput,
		validTypes: []string{"movie"},
		output:     []string{"tt1"},
		pass:       true,
	},
	{
		name:       "combined",
		fail:       false,
		response:   ReferenceOutput,
		validTypes: []string{"movie", "tvSpecial"},
		output:     []string{"tt1", "tt3"},
		pass:       true,
	},
	{
		name:       "invalid",
		fail:       false,
		response:   InvalidOutput,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
	{
		name:       "header missing",
		fail:       false,
		response:   HeaderMissing,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
	{
		name:       "empty",
		fail:       false,
		response:   ``,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
	{
		name:       "error",
		fail:       true,
		response:   ``,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
}

func TestGetByTypes(t *testing.T) {
	handler := Handler{}
	s := httptest.NewServer(http.HandlerFunc(handler.Handle))

	c := watchlist.Client{
		HTTPClient: http.DefaultClient,
		URL:        s.URL,
	}

	for _, test := range GetTestCases {

		handler.Fail = test.fail
		handler.Response = test.response

		entries, err := c.GetByTypes(test.validTypes...)

		if test.pass {
			assert.NoError(t, err, test.name)
			for _, id := range test.output {
				found := func(list []watchlist.Entry) bool {
					for _, entry := range list {
						if entry.IMDBId == id {
							return true
						}
					}
					return false
				}(entries)
				assert.True(t, found, test.name)
			}
		} else {
			assert.Error(t, err, test.name)
		}
	}

	s.Close()
	_, err := c.GetByTypes("movie", "tvSpecial")
	assert.Error(t, err)
}

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
