package link

import (
	"io"
	"log"
	"strings"

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
	var links []Link
	for _, n := range nodes {
		links = append(links, buildLink(n))
	}
	return links
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c) + " "
	}
	return strings.Join(strings.Fields(ret), " ")
}
