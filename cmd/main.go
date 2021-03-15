package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Parser/parser"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type location struct {
	Value string `xml:"loc"`
}

type urlset struct {
	XMLName xml.Name   `xml:"urlset"`
	XMLns   string     `xml:"xmlns,attr"`
	URLs    []location `xml:"url"`
}

type config struct {
	URL        string
	MaxDepth   int
	OutputFile string
}

func main() {
	cfg := config{}
	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)

	_, links := bfs(cfg.URL, cfg.MaxDepth)
	toXML := urlset{
		XMLns: xmlns,
		URLs:  make([]location, len(links)),
	}

	for i, link := range links {
		toXML.URLs[i] = location{
			Value: link,
		}
	}

	if err := marshallXML(cfg.OutputFile, toXML); err != nil {
		log.Println(err)
	}
}

func getLinks(site string) ([]string, []string) {
	resp, err := http.Get(site)
	if err != nil {
		return []string{}, nil
	}
	defer resp.Body.Close()

	links, emails, err := parser.ParsePage(resp.Body)
	if err != nil {
		return []string{}, nil
	}

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}

	base := baseURL.String()

	return filter(buildLink(base, links), withPrefix(base)), trimPfx(emails)
}

func trimPfx(emails []string) []string {
	if emails == nil {
		return emails
	}

	for idx, email := range emails {
		if email != "" {
			emails[idx] = emails[idx][7:]
		}
	}
	return emails
}

func bfs(base string, maxDepth int) ([]string, []string) {
	seen := make(map[string]struct{})
	seenEmail := make(map[string]struct{})
	q := map[string]struct{}{}
	nq := map[string]struct{}{
		base: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		fmt.Println("depth =", i)
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}

		for link := range q {
			if _, ok := seen[link]; ok {
				continue
			}

			seen[link] = struct{}{}
			links, emails := getLinks(link)
			for _, link := range links {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
			}

			for _, link := range emails {
				if _, ok := seenEmail[link]; !ok {
					fmt.Println(link)
					seenEmail[link] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for link := range seen {
		ret = append(ret, link)
	}

	ret2 := make([]string, 0, len(seenEmail))
	for link := range seenEmail {
		ret2 = append(ret2, link)
	}

	return ret, ret2
}

func marshallXML(file string, data urlset) error {
	output, err := xml.MarshalIndent(data, "  ", "    ")
	if err != nil {
		return err
	}

	output = []byte(xml.Header + string(output))

	return ioutil.WriteFile(file, output, 0644)
}

func filter(links []string, KeepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if KeepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(url string) bool {
		return strings.HasPrefix(url, pfx) && !strings.Contains(url, "zip") && !strings.Contains(url, "pdf") && !strings.Contains(url, "pptx") && !strings.Contains(url, "doc") && !strings.Contains(url, "jpg")
	}
}

func buildLink(base string, links []string) []string {
	var ret []string
	for _, link := range links {
		switch {
		case strings.HasPrefix(link, "/"):
			ret = append(ret, base+link)
		case strings.HasPrefix(link, "http"):
			ret = append(ret, link)
		}
	}

	return ret
}
