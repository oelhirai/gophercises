package link

import (
	"io"
	"log"

	"golang.org/x/net/html"
)

// Link represents a link (<a href="...">) in an HTML
// document.
type Link struct {
	Href string
	Text string
}

// Parse will tak in an HTML and will return a
// slice of links parsed from it.
func ParseLinks(r io.Reader) ([]Link, error) {
	// Parse the html file
	doc, err := html.Parse(r)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}

	// Get all nodes with links
	var links []Link
	extractLinksFromNode(doc, &links)
	return links, nil
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
