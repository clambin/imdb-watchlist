package imdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Client fetches an IMDB watchlist and returns the entries that match a set of types
type Client struct {
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
func (c *Client) ReadByTypes(validTypes ...string) ([]Entry, error) {
	allEntries, err := c.getWatchlist()
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

func (c *Client) getWatchlist() ([]Entry, error) {
	url := "https://www.imdb.com"
	if c.URL != "" {
		url = c.URL
	}

	req, _ := http.NewRequest(http.MethodGet, url+"/list/"+c.ListID+"/export", nil)
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return parseList(bytes.NewBuffer(body))
}

func checkType(entryType string, validTypes ...string) bool {
	for _, validType := range validTypes {
		if entryType == validType {
			return true
		}
	}
	return len(validTypes) == 0
}
