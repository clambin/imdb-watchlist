package imdb

import (
	"errors"
	"net/http"
)

// Fetcher fetches an IMDB watchlist and returns the entries that match a set of types
type Fetcher struct {
	HTTPClient *http.Client
	ListID     string
	URL        string
}

// Entry is an entry in an IMDB watchlist
type Entry struct {
	IMDBId string
	Type   string
	Title  string
}

// ReadByTypes queries an IMDB watchlist and returns the entries that match validTypes. If no validTtypes are provided,
// all watchlist entries are returned.
func (f Fetcher) ReadByTypes(validTypes ...string) ([]Entry, error) {
	allEntries, err := f.getWatchlist()
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, entry := range allEntries {
		if checkType(entry.Type, validTypes...) {
			entries = append(entries, entry)
		}
	}
	return entries, err
}

func (f Fetcher) getWatchlist() ([]Entry, error) {
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

func checkType(entryType string, validTypes ...string) bool {
	for _, validType := range validTypes {
		if entryType == validType {
			return true
		}
	}
	return len(validTypes) == 0
}
