package main

import (
	"flag"
	"fmt"

	"github.com/oelhirai/gophercises/link"
)

func main() {
	filename := readFileName()

	var links []link.Link
	links = link.ParseLinks(filename)

	for _, parsedLink := range links {
		fmt.Printf("%+v\n", parsedLink)
	}
}

func readFileName() string {
	filename := flag.String("file", "res/ex1.html", "the HTML file with links to parse")
	flag.Parse()

	fmt.Printf("Using the HTML in %s.\n", *filename)
	return *filename
}
