package watchlist

import (
	"encoding/csv"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Entry map[string]string

func Get(httpClient *http.Client, listID string, validTypes ...string) (entries []Entry, err error) {
	watchlistURL := "https://www.imdb.com/list/" + listID + "/export"

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	var resp *http.Response
	resp, err = httpClient.Get(watchlistURL)

	if err == nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			err = errors.New(resp.Status)
		}
	}

	if err == nil {
		scanner := csv.NewReader(resp.Body)

		var columns []string
		columns, err = scanner.Read()

		if err == nil {
			log.WithFields(log.Fields{"columns": columns, "count": len(columns)}).Debug("column line read")
			entries, err = parseEntries(scanner, columns, validTypes...)
		}
	}

	return
}

func parseEntries(scanner *csv.Reader, columns []string, validTypes ...string) (entries []Entry, err error) {
	var fields []string
	for err == nil {
		if fields, err = scanner.Read(); err == nil {
			newEntry := parseEntry(fields, columns)

			if checkType(newEntry["Title Type"], validTypes...) {
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

func parseEntry(fields []string, columns []string) (entries map[string]string) {
	if len(fields) != len(columns) {
		panic("unexpected csv behaviour: different number of columns & fields")
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
