package sitemap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/oelhirai/gophercises/link"
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
func BuildSiteMap(hostURL string, depth int) {
	var seenLinks *set
	var nextQueue []string
	var currentQueue []string

	// Build link retriever which resolves
	seenLinks = newSet()
	seenLinks.Add(hostURL)
	getHrefs := getLinksRetrieverClosure(hostURL)
	currentQueue, _ = getHrefs(hostURL)

	// Start exploring url in page...
	for depth > 0 {
		for _, l := range currentQueue {
			if !seenLinks.Contains(l) {
				seenLinks.Add(l)
				linksInPage, _ := getHrefs(l)
				nextQueue = append(nextQueue, linksInPage...)
			}
		}
		currentQueue = nextQueue
		nextQueue = make([]string, 0)
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

func getLinksRetrieverClosure(hostSite string) func(string) ([]string, error) {
	return func(site string) ([]string, error) {
		return hrefs(site, hostSite)
	}
}

func hrefs(site string, hostSite string) ([]string, error) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Printf("Error retrieving site: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	var result []string
	links, _ := link.ParseLinks(resp.Body)
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			result = append(result, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			result = append(result, l.Href)
		}
	}
	return result, nil
}
