package page

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	url, time, h1, h4, text string
}

func (i *Info) String() string {
	return strings.Join([]string{
		i.url,
		i.time,
		i.h1,
		i.h4,
		i.text,
		"--------",
	}, "\n")
}

func getText(doc *goquery.Document, s string) string {
	return strings.TrimSpace(doc.Find(s).First().Text())
}

func getInfo(link *url.URL, doc *goquery.Document) *Info {
	text := getText(doc, ".js-article-body")
	if i := strings.IndexByte(text, '\n'); i != -1 {
		text = text[:i]
	}
	return &Info{
		link.String(),
		getText(doc, "time"),
		getText(doc, "h1"),
		getText(doc, "h4"),
		text,
	}
}

func PageInfo(link *url.URL, doc *goquery.Document) *Info {
	path := link.Path
	if len(path) > 0 {
		path = path[1:]

		if i := strings.IndexByte(path, '/'); i != -1 {
			path = path[:i]
		}

		if path == "artikel" && link.Hostname() == "www.breakit.se" {
			return getInfo(link, doc)
		}
	}
	return nil
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
