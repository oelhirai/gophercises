package cyoa

type Story map[string]Arc

type Arc struct {
	Title   string   `json:"title"`
	Paragraphs   []string   `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	NextArc string `json:"arc"`
}