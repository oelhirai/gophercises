package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/oelhirai/gophercises/link"
)

func main() {
	filename := readFileName()
	// Open the html file
	r, err := os.Open(filename)
	if err != nil {
		log.Printf("%v\n", err)
		panic(err)
	}

	links, err := link.ParseLinks(r)
	if err != nil {
		log.Printf("%v\n", err)
		panic(err)
	}

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
