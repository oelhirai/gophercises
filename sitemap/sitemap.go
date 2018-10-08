package sitemap

import (
	"fmt"
	"net/http"
	"os"

	"github.com/oelhirai/gophercises/link"
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

func (s *set) Remove(value string) {
	delete(s.m, value)
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
			if !seenLinks.Contains(l.Href) {
				seenLinks.Add(l.Href)
				linksInPage, _ := getLinks(l.Href)
				nextQueue = append(nextQueue, linksInPage...)
			}
		}
		currentQueue = nextQueue
		nextQueue = make([]link.Link, 0)
		depth--
	}

	for _, v := range seenLinks.m {
		fmt.Println(v)
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

	return links, nil
}
