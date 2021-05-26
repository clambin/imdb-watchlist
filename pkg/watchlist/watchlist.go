package watchlist

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Entry struct {
	IMDBId string
	Type   string
	Title  string
}

func Get(httpClient *http.Client, listID string, validTypes ...string) (entries []Entry, err error) {
	var body io.ReadCloser

	body, err = getWatchlist(httpClient, "https://www.imdb.com/list/"+listID+"/export")

	if err == nil {
		reader := csv.NewReader(body)

		var columns []string
		columns, err = reader.Read()

		var indices map[string]int
		if err == nil {
			indices, err = parseColumns(columns)
		}

		if err == nil {
			log.WithFields(log.Fields{"columns": indices, "count": len(indices)}).Debug("column line read")
			entries, err = parseEntries(reader, indices, validTypes...)
		}

		_ = body.Close()
	}

	return
}

func getWatchlist(httpClient *http.Client, watchlistURL string) (body io.ReadCloser, err error) {
	var resp *http.Response
	resp, err = httpClient.Get(watchlistURL)

	if err == nil {
		body = resp.Body

		if resp.StatusCode != http.StatusOK {
			_ = body.Close()
			err = errors.New(resp.Status)
		}
	}

	return body, err
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
			err = fmt.Errorf("mandatory field missing")
		}
	}

	return
}

func parseEntries(reader *csv.Reader, indices map[string]int, validTypes ...string) (entries []Entry, err error) {
	var fields []string
	for err == nil {
		if fields, err = reader.Read(); err == nil {

			newEntry := Entry{
				IMDBId: fields[indices["Const"]],
				Title:  fields[indices["Title"]],
				Type:   fields[indices["Title Type"]],
			}

			if checkType(newEntry.Type, validTypes...) {
				entries = append(entries, newEntry)
				log.WithFields(log.Fields{"entries": entries, "count": len(fields)}).Debug("entry found")
			}
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
