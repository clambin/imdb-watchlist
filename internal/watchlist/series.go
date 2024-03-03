package watchlist

import (
	"github.com/clambin/imdb-watchlist/pkg/imdb"
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
	entries = Unique(entries, func(v imdb.Entry) string {
		return v.IMDBId
	})
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
