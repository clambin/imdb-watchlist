package imdb

import "log/slog"

// Entry is an entry in an IMDB watchlist
type Entry struct {
	IMDBId string
	Type   EntryType
	Title  string
}

type EntryType string

const (
	Movie        EntryType = "movie"
	TVSeries     EntryType = "tvSeries"
	TVSpecial    EntryType = "tvSpecial"
	TVMiniSeries EntryType = "tvMiniSeries"
)

var _ slog.LogValuer = Watchlist{}

type Watchlist []Entry

func (w Watchlist) Filter(mediaType ...EntryType) Watchlist {
	filtered := make(Watchlist, 0, len(w))
	for _, entry := range w {
		for _, entryType := range mediaType {
			if entry.Type == entryType {
				filtered = append(filtered, entry)
				break
			}
		}
	}
	return filtered
}

func (w Watchlist) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, len(w))
	for _, entry := range w {
		attrs = append(attrs, slog.Group(entry.IMDBId,
			slog.String("title", entry.Title),
			slog.String("type", string(entry.Type)),
		))
	}

	return slog.GroupValue(attrs...)
}
