package sitemap

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/oelhirai/gophercises/link"
	// "github.com/oelhirai/gophercises/link"
)

type set struct {
	m map[string]struct{}
}

func newSet() *set {
	s := &set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *set) Add(value string) {
	var exists = struct{}{}
	s.m[value] = exists
}

func (s *set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func BuildSiteMap(url string, depth int) {
	var seenLinks *set
	var nextQueue []link.Link
	var currentQueue []link.Link

	seenLinks = newSet()
	seenLinks.Add(url)
	currentQueue, err := getLinks(url)
	if err != nil {
		fmt.Printf("Error retrieving site: %s", err)
		os.Exit(1)
	}

	// Start exploring url in page...
	for depth > 0 {
		for _, l := range currentQueue {
			if canCheckURL(l.Href, seenLinks) {
				seenLinks.Add(l.Href)
				linksInPage, _ := getLinks(l.Href)
				nextQueue = append(nextQueue, linksInPage...)
			}
		}
		currentQueue = nextQueue
		nextQueue = make([]link.Link, 0)
		depth--
	}

	for k := range seenLinks.m {
		fmt.Println(k)
	}
}

func getLinks(url string) ([]link.Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error retrieving site: %s", err)
		return nil, err
	}

	links, err := link.ParseLinks(resp.Body)
	if err != nil {
		fmt.Printf("Error retrieving site: %s", err)
		return nil, err
	}

	for _, link := range links {
		fullLink, err := getFullyQualifiedURL(link.Href, url)
		if err != nil {
			fmt.Printf("Error retrieving site: %s", err)
			return nil, err
		}
		link.Href = fullLink.String()
	}

	return links, nil
}

func canCheckURL(curURL string, seenLinks *set) bool {
	if seenLinks.Contains(curURL) {
		return false
	}
	currentURL, _ := url.Parse(curURL)
	if currentURL.Scheme != "https" {
		return false
	}
	return true
}

func getFullyQualifiedURL(currentSite string, referenceSite string) (*url.URL, error) {
	currentURL, err := url.Parse(currentSite)
	if err != nil {
		return nil, err
	}
	if currentURL.Host != "" {
		return currentURL, nil
	}
	referenceURL, err := url.Parse(referenceSite)
	if err != nil {
		return nil, err
	}
	return referenceURL.ResolveReference(currentURL), nil
}
