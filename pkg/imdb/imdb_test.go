package imdb_test

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetByTypes(t *testing.T) {
	tests := []struct {
		name       string
		validTypes []string
		fail       bool
		response   string
		pass       bool
		output     []imdb.Entry
	}{
		{
			name:       "tvSeries",
			validTypes: []string{"tvSeries"},
			fail:       false,
			response:   ReferenceOutput,
			pass:       true,
			output:     []imdb.Entry{{IMDBId: "tt2", Type: "tvSeries", Title: "A TV Series"}},
		},
		{
			name:       "movie",
			validTypes: []string{"movie"},
			fail:       false,
			response:   ReferenceOutput,
			pass:       true,
			output:     []imdb.Entry{{IMDBId: "tt1", Type: "movie", Title: "A Movie"}},
		},
		{
			name:       "combined",
			validTypes: []string{"movie", "tvSpecial"},
			fail:       false,
			response:   ReferenceOutput,
			pass:       true,
			output:     []imdb.Entry{{IMDBId: "tt1", Type: "movie", Title: "A Movie"}, {IMDBId: "tt3", Type: "tvSpecial", Title: "A TV Special"}},
		},
		{
			name:       "error",
			validTypes: []string{"movie", "tvSpecial"},
			fail:       true,
			pass:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Handler{
				Fail:     tt.fail,
				Response: tt.response,
			}
			s := httptest.NewServer(http.HandlerFunc(handler.Handle))
			defer s.Close()

			c := imdb.Client{HTTPClient: http.DefaultClient, URL: s.URL}

			entries, err := c.ReadByTypes(tt.validTypes...)

			if !tt.pass {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err, tt.name)
			assert.Equal(t, tt.output, entries)
		})
	}
}

func TestClient_ReadByTypes_Error(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(nil))
	s.Close()
	c := imdb.Client{HTTPClient: http.DefaultClient, URL: s.URL}
	_, err := c.ReadByTypes()
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
