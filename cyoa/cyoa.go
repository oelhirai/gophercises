package cyoa

import (
	"io"
	"encoding/json"
	"net/http"
	"strings"
	"html/template"
)

func init(){
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf8">
		<title>Create Your Own Adventure</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
			<p>{{.}}</p>
		{{end}}
		{{ if .Options}}
			<ul>
			{{range .Options}}
				<li><a href="{{.NextArc}}">{{.Text}}</a></li>
			{{end}}
			</ul>
		{{ else }}
			<a href="intro">Restart Game</a>
		{{end}}
	</body>
</html>
`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	if story, ok := h.s[path]; ok {
		tpl.Execute(w, story)
	} else {
		tpl.Execute(w, h.s["intro"])
	}
}



func JsonStory(r io.Reader) (Story, error){
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

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