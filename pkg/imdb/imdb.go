package imdb

import (
	"errors"
	"net/http"
)

// WatchlistFetcher fetches an IMDB watchlist and returns the entries that match a set of types.
type WatchlistFetcher struct {
	HTTPClient *http.Client
	URL        string
}

// GetWatchlist queries an IMDB watchlist.
func (f WatchlistFetcher) GetWatchlist(listID string) (Watchlist, error) {
	url := "https://www.imdb.com"
	if f.URL != "" {
		url = f.URL
	}

	req, _ := http.NewRequest(http.MethodGet, url+"/list/"+listID+"/export", nil)
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
