package link

import (
	"log"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseLinks(file string) []Link {
	// Open the html file
	r, err := os.Open(file)
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	// Parse the html file
	doc, err := html.Parse(r)
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	// Get all nodes with links
	var links []Link
	extractLinksFromNode(doc, &links)
	return links
}

func extractLinksFromNode(n *html.Node, links *[]Link) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "a" {
			attr := c.Attr[0]
			*links = append(*links, Link{Href: attr.Val, Text: c.FirstChild.Data})
			continue
		} else {
			extractLinksFromNode(c, links)
		}
	}
}
