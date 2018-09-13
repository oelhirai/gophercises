package main

import (
	"fmt"
	"os"
	"encoding/json"
	"html/template"
	"flag"
	"net/http"
	"strings"
	"github.com/oelhirai/gophercises/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	storyMap := parseJson(*filename)
	handler := createStoryHandler(storyMap)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func parseJson(name string) cyoa.Story {
	jsonFile, err := os.Open(name)
	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(jsonFile)
	var story cyoa.Story
	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	return story
}

func createStoryHandler(storyMap cyoa.Story) http.Handler {
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
				{{range .Paragraphs}}
				  {{.}}
				{{end}}
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