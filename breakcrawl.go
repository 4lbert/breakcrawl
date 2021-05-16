package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func printArticle(link string, doc *goquery.Document) {
	fmt.Println(link)
	time := doc.Find("time").First().Text()
	fmt.Println(strings.TrimSpace(time))
	h1 := doc.Find("h1").First().Text()
	fmt.Println(strings.TrimSpace(h1))
	h4 := doc.Find("h4").First().Text()
	fmt.Println(strings.TrimSpace(h4))
	text := doc.Find(".js-article-body").First().Text()
	text = strings.TrimSpace(text)
	paragraphs := strings.Split(text, "\n")
	if len(paragraphs) > 0 {
		text = paragraphs[0]
	}
	fmt.Println(text)
	fmt.Println("--------")
}

var visited = make(map[string]bool)

func visit(link *url.URL, depth uint, max uint) {
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

	path := strings.Split(link.Path, "/")
	if len(path) > 1 && path[1] == "artikel" {
		printArticle(linkStr, doc)
	}

	if depth < max {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists == false {
				return
			}
			next, err := link.Parse(href)
			if err != nil {
				return
			}

			visit(next, depth+1, max)
		})
	}
}

func main() {
	depth := flag.Uint("d", 1, "the depth of the link traversal")
	flag.Parse()

	link, err := url.Parse("https://www.breakit.se")
	if err != nil {
		log.Fatal(err)
	}
	visit(link, 0, *depth)
}
