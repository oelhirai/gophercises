package main

import (
	"fmt"
	"os"
	"flag"
	"net/http"
	"log"
	"github.com/oelhirai/gophercises/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	storyMap := parseJson(*filename)
	handler := cyoa.NewHandler(storyMap, cyoa.WithTemplate(nil))

	fmt.Printf("Starting the server on %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}

func parseJson(name string) cyoa.Story {
	jsonFile, err := os.Open(name)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(jsonFile)
	if err != nil {
		panic(err)
	}

	return story
}