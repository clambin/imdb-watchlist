package imdb

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
