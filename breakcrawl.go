package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/4lbert/breakcrawl/page"
	"github.com/PuerkitoBio/goquery"
)

var maxDepth uint = 1

type Visit struct {
	url   *url.URL
	depth uint
}

func visit(link *url.URL, depth uint, infoOut chan *page.Info, nextOut chan Visit) {
	res, err := http.Get(link.String())
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	infoOut <- page.PageInfo(link, doc)

	if depth < maxDepth {
		page.ForEachLink(link, doc, func(next *url.URL) {
			nextOut <- Visit{next, depth + 1}
		})
	}
}

func main() {
	flag.UintVar(&maxDepth, "d", 1, "maximum depth of the link traversal")
	maxConcurrent := flag.Uint("p", 8, "maximum number of concurrent outstanding requests")
	flag.Parse()

	fmt.Println(*maxConcurrent)

	link, err := url.Parse("https://www.breakit.se")
	if err != nil {
		log.Fatal(err)
	}

	visited := make(map[string]bool)
	info := make(chan *page.Info)
	next := make(chan Visit)

	go visit(link, 0, info, next)

	for {
		select {
		case i := <-info:
			if i != nil {
				fmt.Println(i)
			}
		case n := <-next:
			s := n.url.String()
			if _, ok := visited[s]; !ok {
				go visit(n.url, n.depth, info, next)
				visited[s] = true
			}
		}
	}
}
