package watchlist

import (
	"cmp"
	"github.com/clambin/imdb-watchlist/pkg/imdb"
	"slices"
)

func (s *Server) getSeries() ([]Entry, error) {
	var entries []imdb.Entry
	for _, r := range s.readers {
		newEntries, err := r.GetWatchlist("tvSeries", "tvMiniSeries")
		if err != nil {
			return nil, err
		}
		entries = append(entries, newEntries...)
	}
	return s.buildSeriesResponse(unique(entries)), nil
}

func unique(input []imdb.Entry) []imdb.Entry {
	slices.SortFunc(input, func(a, b imdb.Entry) int {
		return cmp.Compare(a.IMDBId, b.IMDBId)
	})
	var lastID string
	entries := make([]imdb.Entry, 0, len(input))
	for _, e := range input {
		if e.IMDBId == lastID {
			continue
		}
		entries = append(entries, e)
		lastID = e.IMDBId
	}
	return entries
}

// Entry represents an entry in the IMDB watchlist
type Entry struct {
	Title  string `json:"title"`
	IMDBId string `json:"imdbId"`
}

func (s *Server) buildSeriesResponse(imdbEntries []imdb.Entry) []Entry {
	entries := make([]Entry, len(imdbEntries))

	for i := range imdbEntries {
		entries[i] = Entry{
			Title:  imdbEntries[i].Title,
			IMDBId: imdbEntries[i].IMDBId,
		}

		s.logger.Debug("imdb watchlist entry found",
			"title", imdbEntries[i].Title,
			"imdbId", imdbEntries[i].IMDBId,
			"count", i+1,
		)
	}

	return entries
}
