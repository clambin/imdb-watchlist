package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/watchlist"
	"github.com/clambin/imdb-watchlist/watchlist/server"
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
		response:   server.ReferenceOutput,
		validTypes: []string{"tvSeries"},
		output:     []string{"tt2"},
		pass:       true,
	},
	{
		name:       "movie",
		fail:       false,
		response:   server.ReferenceOutput,
		validTypes: []string{"movie"},
		output:     []string{"tt1"},
		pass:       true,
	},
	{
		name:       "combined",
		fail:       false,
		response:   server.ReferenceOutput,
		validTypes: []string{"movie", "tvSpecial"},
		output:     []string{"tt1", "tt3"},
		pass:       true,
	},
	{
		name:       "invalid",
		fail:       false,
		response:   server.InvalidOutput,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
	{
		name:       "header missing",
		fail:       false,
		response:   server.HeaderMissing,
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
	handler := server.Handler{}
	s := httptest.NewServer(http.HandlerFunc(handler.Handle))

	client := watchlist.Client{URL: s.URL}

	for _, test := range GetTestCases {

		handler.Fail = test.fail
		handler.Response = test.response

		entries, err := client.GetByTypes(test.validTypes...)

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
	_, err := client.GetByTypes("movie", "tvSpecial")
	assert.Error(t, err)
}
