package page

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getText(doc *goquery.Document, s string) string {
	return strings.TrimSpace(doc.Find(s).First().Text())
}

func printArticle(link *url.URL, doc *goquery.Document) {
	fmt.Println(link.String())
	fmt.Println(getText(doc, "time"))
	fmt.Println(getText(doc, "h1"))
	fmt.Println(getText(doc, "h4"))
	text := getText(doc, ".js-article-body")
	if i := strings.IndexByte(text, '\n'); i != -1 {
		text = text[:i]
	}
	fmt.Println(text)
	fmt.Println("--------")
}

func PrintPage(link *url.URL, doc *goquery.Document) {
	path := link.Path
	if len(path) > 0 {
		path = path[1:]

		if i := strings.IndexByte(path, '/'); i != -1 {
			path = path[:i]
		}

		if path == "artikel" && link.Hostname() == "www.breakit.se" {
			printArticle(link, doc)
		}
	}
}

func ForEachLink(link *url.URL, doc *goquery.Document, fn func(next *url.URL)) {
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists == false {
			return
		}
		next, err := link.Parse(href)
		if err != nil {
			return
		}
		fn(next)
	})
}
