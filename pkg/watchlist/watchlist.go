package watchlist

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/clambin/go-metrics/client"
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
	client.Caller
	ListID string
	URL    string
}

var _ Reader = &Client{}

// Entry is an entry in an IMDB watchlist
type Entry struct {
	IMDBId string
	Type   string
	Title  string
}

// GetAll queries an IMDB watchlist and returns all entries
func (client *Client) GetAll() (entries []Entry, err error) {
	url := "https://www.imdb.com"
	if client.URL != "" {
		url = client.URL
	}

	req, _ := http.NewRequest(http.MethodGet, url+"/list/"+client.ListID+"/export", nil)

	var resp *http.Response
	resp, err = client.Caller.Do(req)

	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	return parseList(resp.Body)
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

func parseList(body io.ReadCloser) (entries []Entry, err error) {
	reader := csv.NewReader(body)

	var columns []string
	columns, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	var indices map[string]int
	indices, err = parseColumns(columns)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	log.WithFields(log.Fields{"columns": indices, "count": len(indices)}).Debug("column line read")
	return parseEntries(reader, indices)
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
		if !found {
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
