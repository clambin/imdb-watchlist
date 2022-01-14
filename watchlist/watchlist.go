package watchlist

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// Reader interface fetches an IMDB watchlist and returns the entries that match validTypes
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

// GetAll queries an IDMB watchlist and returns all entries
func (client *Client) GetAll() (entries []Entry, err error) {
	var body io.ReadCloser
	body, err = client.getWatchlist(client.ListID)
	if err != nil {
		return nil, fmt.Errorf("get failed: %w", err)
	}

	defer func() {
		_ = body.Close()
	}()

	reader := csv.NewReader(body)

	var columns []string
	columns, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	var indices map[string]int
	indices, err = parseColumns(columns)

	log.WithFields(log.Fields{"columns": indices, "count": len(indices)}).Debug("column line read")
	entries, err = parseEntries(reader, indices)

	if err != nil {
		err = fmt.Errorf("parse failed: %w", err)
	}

	return
}

// GetByTypes queries an IMDB watchlist and returns the entries that match validTypes
func (client *Client) GetByTypes(validTypes ...string) (entries []Entry, err error) {
	var allEntries []Entry
	allEntries, err = client.GetAll()

	if err != nil {
		return
	}

	for _, entry := range allEntries {
		if checkType(entry.Type, validTypes...) {
			entries = append(entries, entry)
		}
	}

	return
}

func (client *Client) getWatchlist(listID string) (body io.ReadCloser, err error) {
	if client.HTTPClient == nil {
		client.HTTPClient = http.DefaultClient
	}
	url := "https://www.imdb.com"
	if client.URL != "" {
		url = client.URL
	}

	watchListURL := url + "/list/" + listID + "/export"

	var resp *http.Response
	resp, err = client.HTTPClient.Get(watchListURL)

	if err != nil {
		return
	}

	body = resp.Body

	if resp.StatusCode != http.StatusOK {
		_ = body.Close()
		err = errors.New(resp.Status)
	}

	return
}

func parseColumns(columns []string) (indices map[string]int, err error) {
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
		if found == false {
			log.WithField("column", column).Error("mandatory field missing")
			return nil, fmt.Errorf("mandatory field '%s' missing", column)
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
