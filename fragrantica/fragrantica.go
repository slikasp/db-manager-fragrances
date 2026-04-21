package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

// TODOS

type Scraper struct {
	client  tls_client.HttpClient
	headers http.Header
}

type Fragrance struct {
	ID            int64
	Url           string
	Name          string
	Brand         string
	Country       string
	Gender        string
	RatingValue   string
	RatingCount   int32
	Year          int32
	TopNotes      string
	MiddleNotes   string
	BaseNotes     string
	Perfumer1     string
	Perfumer2     string
	Accord1       string
	Accord2       string
	Accord3       string
	Accord4       string
	Accord5       string
	FragranticaID int32
}

// HTTP client that looks like a real browser to avoid being blocked
func NewScraper() (*Scraper, error) {
	jar := tls_client.NewCookieJar()

	client, err := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithCookieJar(jar),
		tls_client.WithRandomTLSExtensionOrder(),
	)
	if err != nil {
		return nil, err
	}

	headers := http.Header{
		"accept": {
			"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		},
		"accept-encoding": {"gzip, deflate, br"},
		"accept-language": {"en-US,en;q=0.9"},
		"cache-control":   {"max-age=0"},
		"sec-ch-ua": {
			`"Chromium";v="120", "Not(A:Brand";v="24", "Google Chrome";v="120"`,
		},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-user":            {"?1"},
		"upgrade-insecure-requests": {"1"},
		"user-agent": {
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"sec-fetch-user",
			"upgrade-insecure-requests",
			"user-agent",
		},
	}

	return &Scraper{client: client, headers: headers}, nil
}

// Get html from url
func (s *Scraper) GetPageBody(url string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = s.headers

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// Get details of newly found fragrance by parsing html page
// TODO: try using ML and get the details from the cards...
func ParsePage(doc *goquery.Document) {

	// accords := getAccords(doc)

}

func getAccords(doc *goquery.Document) []string {
	var results []string

	doc.Find("body > main > div.flex.flex-col.w-full.max-w-\\[280px\\].md\\:max-w-\\[320px\\]").Each(func(i int, s *goquery.Selection) {
		s.Find("> div > div > span.truncate").Each(func(j int, item *goquery.Selection) {
			results = append(results, item.Text())
		})
	})

	return results
}

// func - (to be run periodically (once an hour or so) take one existing frag and update the card and public score (anything else?) if changed
