package sitemap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/oelhirai/gophercises/link"
	// "github.com/oelhirai/gophercises/link"
)

// PageURL is a struct simply holding a url
type PageURL struct {
	URL string `xml:"loc"`
}

// URLSet is a struct used in converting data to XML
type URLSet struct {
	XMLName xml.Name  `xml:"urlset"`
	Urls    []PageURL `xml:"url"`
}

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

// BuildSiteMap extracts all links from the given host.
// the depth is the maximum number of links to follow when building the sitemap
func BuildSiteMap(site string, depth int) {
	var seenLinks *set
	var nextQueue []link.Link
	var currentQueue []link.Link

	seenLinks = newSet()
	seenLinks.Add(site)
	currentQueue, err := getLinks(site)
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

	sitemap := &URLSet{}
	for k := range seenLinks.m {
		sitemap.Urls = append(sitemap.Urls, PageURL{k})
	}

	output, err := xml.MarshalIndent(sitemap, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)
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
		fullLink, err := getAbsoluteURL(link.Href, url)
		if err != nil {
			fmt.Printf("Error retrieving site: %s", err)
			return nil, err
		}
		link.Href = fullLink.String()
	}

	return links, nil
}

func canCheckURL(curSite string, seenLinks *set) bool {
	if seenLinks.Contains(curSite) {
		return false
	}
	curURL, _ := url.Parse(curSite)
	if curURL.Scheme != "https" {
		return false
	}
	return true
}

func getAbsoluteURL(currentSite string, referenceSite string) (*url.URL, error) {
	currentURL, err := url.Parse(currentSite)
	if err != nil {
		return nil, err
	}
	if currentURL.IsAbs() {
		return currentURL, nil
	}

	// currentUrl is relative, resolve with host of reference site
	referenceURL, err := url.Parse(referenceSite)
	if err != nil {
		return nil, err
	}
	return referenceURL.ResolveReference(currentURL), nil
}
