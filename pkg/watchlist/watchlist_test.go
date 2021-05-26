package watchlist_test

import (
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/imdb-watchlist/pkg/watchlist"
	"github.com/clambin/imdb-watchlist/pkg/watchlist/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCase struct {
	name       string
	input      string
	validTypes []string
	output     []string
	pass       bool
}

var GetTestCases = []TestCase{
	{
		name:       "tvSeries",
		input:      mock.ReferenceOutput,
		validTypes: []string{"tvSeries"},
		output:     []string{"tt2"},
		pass:       true,
	},
	{
		name:       "movie",
		input:      mock.ReferenceOutput,
		validTypes: []string{"movie"},
		output:     []string{"tt1"},
		pass:       true,
	},
	{
		name:       "combined",
		input:      mock.ReferenceOutput,
		validTypes: []string{"movie", "tvSpecial"},
		output:     []string{"tt1", "tt3"},
		pass:       true,
	},
	{
		name:       "invalid",
		input:      mock.InvalidOutput,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
	{
		name:       "empty",
		input:      mock.EmptyOutput,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
}

func TestGet(t *testing.T) {
	server := httpstub.NewTestClient(mock.Serve)

	for _, test := range GetTestCases {

		mock.ServerOutput = test.input

		entries, err := watchlist.Get(server, "1", test.validTypes...)

		if test.pass {
			assert.NoError(t, err, test.name)
			for _, id := range test.output {
				found := func(list []watchlist.Entry) bool {
					for _, entry := range list {
						if value, ok := entry["Const"]; ok {
							if value == id {
								return true
							}
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

}
