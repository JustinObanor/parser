package parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func ParsePage(r io.Reader) ([]string, []string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, nil, err
	}

	links, emails := getLinks(doc)
	return links, emails, nil
}

func getLinks(n *html.Node) ([]string, []string) {
	email := ""
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				if strings.HasPrefix(attr.Val, "mailto") {
					email = attr.Val
				}
				return []string{attr.Val}, []string{email}
			}
		}
	}

	var ret []string
	var ret2 []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links, emails := getLinks(c)
		ret = append(ret, links...)
		if emails != nil {
			ret2 = append(ret2, emails...)
		}
	}

	return ret, ret2
}
