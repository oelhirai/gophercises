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
	linkNodes := linkNodes(doc)
	links := convertNodesToLinks(linkNodes)
	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	var nodes []*html.Node
	if n.Type == html.ElementNode && n.Data == "a" {
		nodes = append(nodes, n)
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			nodes = append(nodes, linkNodes(c)...)
		}
	}
	return nodes
}

func convertNodesToLinks(nodes []*html.Node) []Link {
	links := make([]Link, len(nodes))
	for _, n := range nodes {
		attr := n.Attr[0]
		links = append(links, Link{Href: attr.Val, Text: n.FirstChild.Data})
	}
	return links
}
