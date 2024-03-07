package imdb

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
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

func getColumnIndices(columns []string) (map[string]int, error) {
	var mandatory = map[string]struct{}{
		"Const":      {},
		"Title":      {},
		"Title Type": {},
	}

	indices := make(map[string]int, len(columns))
	for index, column := range columns {
		indices[column] = index
		delete(mandatory, column)
	}

	if len(mandatory) > 0 {
		missing := make([]string, 0, len(mandatory))
		for column := range mandatory {
			missing = append(missing, column)
		}
		return nil, fmt.Errorf("imdb: mandatory fields missing: %s", strings.Join(missing, ","))
	}

	return indices, nil
}

func parseEntries(reader *csv.Reader, indices map[string]int) (entries []Entry, err error) {
	var fields []string
	for err == nil {
		if fields, err = reader.Read(); err == nil {
			entries = append(entries, Entry{
				IMDBId: fields[indices["Const"]],
				Title:  fields[indices["Title"]],
				Type:   EntryType(fields[indices["Title Type"]]),
			})
		}
	}

	if err == io.EOF {
		err = nil
	}
	return
}
