package watchlist

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Reader interface fetches an IMDB watchlist and returns the entries that match validTypes
//
//go:generate mockery --name Reader
type Reader interface {
	GetAll() (entries []Entry, err error)
	GetByTypes(validTypes ...string) (entries []Entry, err error)
}

// Client fetches an IMDB watchlist and returns the entries that match a set of types
type Client struct {
	HTTPClient *http.Client
	ListID     string
	URL        string
}

var _ Reader = &Client{}

// Entry is an entry in an IMDB watchlist
type Entry struct {
	IMDBId string
	Type   string
	Title  string
}

// GetAll queries an IMDB watchlist and returns all entries
func (client *Client) GetAll() ([]Entry, error) {
	url := "https://www.imdb.com"
	if client.URL != "" {
		url = client.URL
	}

	req, _ := http.NewRequest(http.MethodGet, url+"/list/"+client.ListID+"/export", nil)
	resp, err := client.HTTPClient.Do(req)

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

// GetByTypes queries an IMDB watchlist and returns the entries that match validTypes
func (client *Client) GetByTypes(validTypes ...string) ([]Entry, error) {
	allEntries, err := client.GetAll()
	var entries []Entry
	if err == nil {
		for _, entry := range allEntries {
			if checkType(entry.Type, validTypes...) {
				entries = append(entries, entry)
			}
		}
	}
	return entries, err
}

func parseList(body io.ReadCloser) ([]Entry, error) {
	reader := csv.NewReader(body)

	columns, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("csv: %w", err)
	}

	var indices map[string]int
	indices, err = getColumnIndices(columns)
	if err != nil {
		return nil, err
	}

	return parseEntries(reader, indices)
}

func getColumnIndices(columns []string) (indices map[string]int, err error) {
	var mandatory = map[string]bool{
		"Const":      false,
		"Title":      false,
		"Title Type": false,
	}

	indices = make(map[string]int)
	for index, column := range columns {
		indices[column] = index
		mandatory[column] = true
	}

	for column, found := range mandatory {
		if !found {
			err = fmt.Errorf("watchlist: mandatory field '%s' missing", column)
			break
		}
	}

	return
}

func parseEntries(reader *csv.Reader, indices map[string]int) (entries []Entry, err error) {
	var fields []string
	for err == nil {
		if fields, err = reader.Read(); err == nil {
			entries = append(entries, Entry{
				IMDBId: fields[indices["Const"]],
				Title:  fields[indices["Title"]],
				Type:   fields[indices["Title Type"]],
			})
		}
	}

	if err == io.EOF {
		err = nil
	}
	return
}

func checkType(entryType string, validTypes ...string) bool {
	for _, validType := range validTypes {
		if entryType == validType {
			return true
		}
	}
	return len(validTypes) == 0
}
