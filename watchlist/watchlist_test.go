package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/watchlist"
	"github.com/clambin/imdb-watchlist/watchlist/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	name       string
	fail       bool
	invalid    bool
	validTypes []string
	output     []string
	pass       bool
}

var GetTestCases = []TestCase{
	{
		name:       "tvSeries",
		fail:       false,
		invalid:    false,
		validTypes: []string{"tvSeries"},
		output:     []string{"tt2"},
		pass:       true,
	},
	{
		name:       "movie",
		fail:       false,
		invalid:    false,
		validTypes: []string{"movie"},
		output:     []string{"tt1"},
		pass:       true,
	},
	{
		name:       "combined",
		fail:       false,
		invalid:    false,
		validTypes: []string{"movie", "tvSpecial"},
		output:     []string{"tt1", "tt3"},
		pass:       true,
	},
	{
		name:       "invalid",
		fail:       false,
		invalid:    true,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
	{
		name:       "empty",
		fail:       true,
		invalid:    false,
		validTypes: []string{"movie", "tvSpecial"},
		pass:       false,
	},
}

func TestGet(t *testing.T) {
	handler := mock.Handler{}
	server := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer server.Close()

	client := watchlist.Client{URL: server.URL}

	for _, test := range GetTestCases {

		handler.Fail = test.fail
		handler.Invalid = test.invalid

		entries, err := client.Watchlist("1", test.validTypes...)

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

}
