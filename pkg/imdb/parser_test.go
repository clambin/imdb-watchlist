package imdb

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pass     bool
		expected []Entry
	}{
		{
			name: "valid",
			input: `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
4,tt4,,,,A TV miniseries,,tvMiniSeries,,,,,,,
`,
			pass: true,
			expected: []Entry{
				{IMDBId: "tt1", Type: "movie", Title: "A Movie"},
				{IMDBId: "tt2", Type: "tvSeries", Title: "A TV Series"},
				{IMDBId: "tt3", Type: "tvSpecial", Title: "A TV Special"},
				{IMDBId: "tt4", Type: "tvMiniSeries", Title: "A TV miniseries"},
			},
		},
		{
			name: "empty",
			input: `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
`,
			pass:     true,
			expected: nil,
		},
		{
			name:  "no input",
			input: ``,
			pass:  false,
		},
		{
			name: "missing header",
			input: `1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
4,tt4,,,,A TV miniseries,,tvMiniSeries,,,,,,,
`,
			pass: false,
		},
		{
			name: "invalid",
			input: `Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,tt1,,,,A Movie,,movie,,,,,,,
2,tt2,,,,A TV Series,,tvSeries,,,,,,,,
3,tt3,,,,A TV Special,,tvSpecial,,,,,,,
`,
			pass: false,
		},
		{
			name: "invalid header",
			input: `Position,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
1,,,,A Movie,,movie,,,,,,,
2,,,,A TV Series,,tvSeries,,,,,,,,
3,,,,A TV Special,,tvSpecial,,,,,,,
`,
			pass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.input)
			entries, err := parseList(r)

			if !tt.pass {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, entries)
		})
	}
}
