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
	mediaTypes := make(map[EntryType]struct{})
	for i := range mediaType {
		mediaTypes[mediaType[i]] = struct{}{}
	}

	filtered := make(Watchlist, 0, len(w))
	for _, entry := range w {
		if _, ok := mediaTypes[entry.Type]; ok {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
