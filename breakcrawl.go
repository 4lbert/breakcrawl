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
	url   url.URL
	depth uint
}

func visit(link url.URL, depth uint, infoOut chan *page.Info, nextOut chan Visit, done chan bool) {
	res, err := http.Get(link.String())
	if err != nil {
		done <- true
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		done <- true
		return
	}

	infoOut <- page.PageInfo(&link, doc)

	if depth < maxDepth {
		page.ForEachLink(&link, doc, func(next *url.URL) {
			nextOut <- Visit{*next, depth + 1}
		})
	}

	done <- true
}

func main() {
	flag.UintVar(&maxDepth, "d", 1, "maximum depth of the link traversal")
	maxConcurrent := flag.Uint("p", 8, "maximum number of concurrent outstanding requests")
	flag.Parse()

	link, err := url.Parse("https://www.breakit.se")
	if err != nil {
		log.Fatal(err)
	}

	visited := make(map[string]bool)
	info := make(chan *page.Info)
	next := make(chan Visit)
	done := make(chan bool)

	go visit(*link, 0, info, next, done)

	var concurrent uint = 1
	var queue []Visit

	popQueue := func() {
		if len(queue) > 0 && concurrent <= *maxConcurrent {
			v := queue[0]
			go visit(v.url, v.depth, info, next, done)
			queue = queue[1:]
			concurrent += 1
		}
	}

	for {
		select {
		case i := <-info:
			if i != nil {
				fmt.Println(i)
			}
		case n := <-next:
			s := n.url.String()
			if _, ok := visited[s]; !ok {
				queue = append(queue, n)
				visited[s] = true
				popQueue()
			}
		case <-done:
			concurrent -= 1
			popQueue()
			if concurrent == 0 && len(queue) == 0 {
				return
			}
		}
	}
}
