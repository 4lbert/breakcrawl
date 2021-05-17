package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/4lbert/breakcrawl/page"
	"github.com/PuerkitoBio/goquery"
)

var maxDepth uint = 1

var visited = make(map[string]bool)

func visit(link *url.URL, depth uint) {
	linkStr := link.String()
	if _, ok := visited[linkStr]; ok {
		return
	}
	visited[linkStr] = true

	res, err := http.Get(linkStr)
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	page.PrintPage(link, doc)

	if depth < maxDepth {
		page.ForEachLink(link, doc, func(next *url.URL) {
			visit(next, depth+1)
		})
	}
}

func main() {
	flag.UintVar(&maxDepth, "d", 1, "the depth of the link traversal")
	flag.Parse()

	link, err := url.Parse("https://www.breakit.se")
	if err != nil {
		log.Fatal(err)
	}
	visit(link, 0)
}
