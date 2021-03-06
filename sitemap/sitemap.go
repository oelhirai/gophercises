package sitemap

import (
	"encoding/xml"
	"fmt"
	"io"
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
	_, ok := s.m[value]
	return ok
}

// BuildSiteMap extracts all links from the given host.
// the depth is the maximum number of links to follow when building the sitemap
func BuildSiteMap(hostURL string, depth int) {
	seenLinks := newSet()
	nextQueue := []string{hostURL}
	var currentQueue []string

	for i := 0; i <= depth; i++ {
		currentQueue, nextQueue = nextQueue, make([]string, 0)
		for _, ref := range currentQueue {
			if !seenLinks.Contains(ref) {
				seenLinks.Add(ref)
				nextQueue = append(nextQueue, get(ref)...)
			}
		}
		currentQueue = nextQueue
	}

	sitemap := URLSet{}
	for k := range seenLinks.m {
		sitemap.Urls = append(sitemap.Urls, PageURL{k})
	}

	output, err := xml.MarshalIndent(sitemap, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)
}

func get(site string) []string {
	resp, err := http.Get(site)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	var result []string
	links, _ := link.ParseLinks(r)
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			result = append(result, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			result = append(result, l.Href)
		}
	}
	return result
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, l := range links {
		if keepFn(l) {
			ret = append(ret, l)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(l string) bool {
		return strings.HasPrefix(l, pfx)
	}
}
