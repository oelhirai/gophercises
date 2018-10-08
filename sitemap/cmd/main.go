package main

import (
	"flag"

	"github.com/oelhirai/gophercises/sitemap"
)

func main() {
	// BFS
	// 1. Create a list of seen links
	// 2. Create a queue of links to explore
	// foreach link, if not in seen set
	//   a. Add link to seen set
	//   b. follow page and add all links to queue (there will be repeats, no worries)

	var site = flag.String("site", "https://www.calhoun.io/", "This is the website we're build a sitemap for")
	var depth = flag.Int("depth", 0, "Depth of the sitemap search")
	flag.Parse()

	sitemap.BuildSiteMap(*site, *depth)
}
