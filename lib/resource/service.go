package resource

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/holmes89/magpie/lib"
	v1 "github.com/holmes89/magpie/lib/handlers/rest/v1"
	"golang.org/x/net/html"
)

type service struct {
	v1.SiteService
}

func NewService(svc v1.SiteService) v1.SiteService {
	return &service{
		svc,
	}
}

func (s *service) Create(ctx context.Context, si lib.Site) error {
	site := si
	title, sites, err := analyze(site.URL)
	if err != nil {
		fmt.Printf("Error analyzing page: %s %s", site.URL, err)
	}
	site.Name = title
	site.Meta = make(map[string]interface{})
	site.Meta["links"] = sites
	return s.SiteService.Create(ctx, site)
}

type subsite struct {
	URL   string `json:"url" dynamodbav:"url"`
	Title string `json:"title" dynamodbav:"title"`
}

func analyze(url string) (string, []subsite, error) {
	page, err := parse(url)
	if err != nil {
		fmt.Printf("Error getting page %s %s\n", url, err)
		return "", nil, errors.New("unable to find page")
	}
	title := pageTitle(page)

	subsites := []subsite{}
	links := pageLinks(page)
	for _, link := range links {
		t, err := pageTitleByURL(link)
		if err != nil {
			fmt.Printf("Error getting page title for %s %s\n", link, err)
		}
		subsites = append(subsites, subsite{
			URL:   link,
			Title: t,
		})
	}

	return title, subsites, nil
}

func pageTitleByURL(url string) (string, error) {
	page, err := parse(url)
	if err != nil {
		fmt.Printf("Error getting page %s %s\n", url, err)
		return "", errors.New("unable to find page")
	}
	t := pageTitle(page)
	return t, nil
}

func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

func pageLinks(n *html.Node) []string {
	links := []string{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		}
	}
	return links
}

func parse(url string) (*html.Node, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Cannot get page")
	}
	b, err := html.Parse(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse page")
	}
	return b, err
}
