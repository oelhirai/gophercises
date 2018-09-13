package main

import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"html/template"
	"net/http"
	"strings"
)

const jsonFileName string = "gopher.json"

type Arc struct {
	Title   string   `json:"title"`
	Story   []string   `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	NextArc string `json:"arc"`
}

func main() {
	// 1. Parse the json into a map (key -> story type)
	storyMap := parseJson(jsonFileName)

	// 2. Learn how to use http/template
	handler := createStoryHandler(storyMap)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func parseJson(name string) map[string]Arc {
	jsonFile, _ := os.Open(name)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jsonMap map[string]Arc
	json.Unmarshal(byteValue, &jsonMap)
	return jsonMap
}

func createStoryHandler(storyMap map[string]Arc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimLeft(r.URL.Path, "/")
		t := template.New("Arc Template")
		t, _ = t.Parse(`
		<html>
		<head>
			<title>Create Your Own Adventure!</title>
		</head>
		<body>
			<h2>{{.Title}}</h2>
			<p>
				{{index .Story 0}}
			</p>

			{{ if .Options}}
				<h3>Make a choice...</h3>
				{{range .Options}}
					<a href="{{.NextArc}}">{{.NextArc}}</a> : {{.Text}}
					<br>
				{{end}}
			{{ else }}
				<a href="intro">Restart Game</a>
			{{end}}
		</body>
		</html>
		`)

		if story, ok := storyMap[path]; ok {
			t.Execute(w, story)
		} else {
			t.Execute(w, storyMap["intro"])
		}
	})
}