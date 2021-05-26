package watchlist

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Entry map[string]string

func Get(httpClient *http.Client, listID string, validTypes ...string) (entries []Entry, err error) {
	var body io.ReadCloser

	body, err = getWatchlist(httpClient, "https://www.imdb.com/list/"+listID+"/export")

	if err == nil {
		reader := csv.NewReader(body)

		var columns []string
		columns, err = reader.Read()

		if err == nil {
			log.WithFields(log.Fields{"columns": columns, "count": len(columns)}).Debug("column line read")
			entries, err = parseEntries(reader, columns, validTypes...)
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

func parseEntries(reader *csv.Reader, columns []string, validTypes ...string) (entries []Entry, err error) {
	var fields []string
	for err == nil {
		if fields, err = reader.Read(); err == nil {
			var newEntry map[string]string
			newEntry, err = parseEntry(fields, columns)

			if err == nil && checkType(newEntry["Title Type"], validTypes...) {
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

func parseEntry(fields []string, columns []string) (entries map[string]string, err error) {
	if len(fields) != len(columns) {
		err = fmt.Errorf("unexpected csv behaviour: %d columns but %d fields", len(columns), len(fields))
		return
	}

	entries = make(map[string]string)
	for index, field := range fields {
		entries[columns[index]] = field
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
