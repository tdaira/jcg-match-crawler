package extender

type Deck struct {
	Name     string
	URL      string
	DeckType string
}

type Win struct {
	Name     string
	DeckType string
}

type MatchInfo struct {
	DeckList []Deck
	WinList  []Win
}
