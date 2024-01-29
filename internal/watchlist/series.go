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

func (s *Server) buildSeriesResponse(entries []imdb.Entry) []Entry {
	sonarrEntries := make([]Entry, 0)

	for _, entry := range entries {
		sonarrEntries = append(sonarrEntries, Entry{
			Title:  entry.Title,
			IMDBId: entry.IMDBId,
		})

		s.logger.Debug("imdb watchlist entry found",
			"title", entry.Title,
			"imdbId", entry.IMDBId,
			"count", len(sonarrEntries),
		)
	}

	return sonarrEntries
}
