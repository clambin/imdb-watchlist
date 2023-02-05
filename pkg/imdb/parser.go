package imdb

import (
	"encoding/csv"
	"fmt"
	"io"
)

func parseList(body io.Reader) ([]Entry, error) {
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
			err = fmt.Errorf("imdb: mandatory field '%s' missing", column)
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
