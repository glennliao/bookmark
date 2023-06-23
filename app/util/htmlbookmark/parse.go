package htmlbookmark

import (
	"context"
	"strings"

	"golang.org/x/net/html"
)

type Bookmarks struct {
	Title    string
	URL      string
	IsDir    bool
	Children []*Bookmarks
}

func Parse(ctx context.Context, data string) (*Bookmarks, error) {
	var file = strings.NewReader(data)

	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node, *Bookmarks) *Bookmarks
	f = func(n *html.Node, b *Bookmarks) *Bookmarks {
		if n.Type == html.ElementNode && n.Data == "a" {
			var url, text string
			for _, a := range n.Attr {
				if a.Key == "href" {
					url = a.Val
					break
				}
			}
			if n.FirstChild != nil {
				text = n.FirstChild.Data
			}
			b.Children = append(b.Children, &Bookmarks{URL: url, Title: text})
			return b
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "h3" {
				var text string
				if c.FirstChild != nil {
					text = c.FirstChild.Data
				}
				dir := &Bookmarks{Title: text, IsDir: true}
				b.Children = append(b.Children, dir)
				b = dir
			} else if c.Type != html.TextNode && c.Type != html.DoctypeNode && c.Type != html.CommentNode {
				f(c, b)
			}
		}
		return b
	}
	res := f(doc, &Bookmarks{Title: "Bookmarks", IsDir: true})

	return res, nil
}
