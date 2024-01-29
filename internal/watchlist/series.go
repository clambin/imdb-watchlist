package watchlist

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
)

func (s *Server) getSeries() ([]Entry, error) {
	entries, err := s.reader.GetWatchlist("tvSeries", "tvMiniSeries")
	if err != nil {
		return nil, err
	}

	return s.buildSeriesResponse(entries), nil
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
