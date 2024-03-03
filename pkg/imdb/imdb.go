package imdb

import (
	"errors"
	"github.com/clambin/go-common/set"
	"net/http"
)

// WatchlistFetcher fetches an IMDB watchlist and returns the entries that match a set of types
type WatchlistFetcher struct {
	HTTPClient *http.Client
	ListID     string
	URL        string
}

// Entry is an entry in an IMDB watchlist
type Entry struct {
	IMDBId string
	Type   EntryType
	Title  string
}

type EntryType string

const (
	Movie        EntryType = "movie"
	TVSeries     EntryType = "tvSeries"
	TVSpecial    EntryType = "tvSpecial"
	TVMiniSeries EntryType = "tvMiniSeries"
)

// GetWatchlist queries an IMDB watchlist and returns the entries that match validTypes. If no validTtypes are provided,
// all watchlist entries are returned.
func (f WatchlistFetcher) GetWatchlist(validTypes ...EntryType) ([]Entry, error) {
	allEntries, err := f.getWatchlist()
	if err != nil {
		return nil, err
	}

	valid := set.New(validTypes...)

	var entries []Entry
	for _, entry := range allEntries {
		if valid.Contains(entry.Type) {
			entries = append(entries, entry)
		}
	}
	return entries, err
}

func (f WatchlistFetcher) getWatchlist() ([]Entry, error) {
	url := "https://www.imdb.com"
	if f.URL != "" {
		url = f.URL
	}

	req, _ := http.NewRequest(http.MethodGet, url+"/list/"+f.ListID+"/export", nil)
	resp, err := f.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return parseList(resp.Body)
}
