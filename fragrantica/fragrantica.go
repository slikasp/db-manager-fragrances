package fragrantica

import (
	"fmt"
	"strconv"
	"strings"

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

type FragranceParams struct {
	FragranticaID int32
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
	Accord6       string
	Accord7       string
	Accord8       string
	Accord9       string
	Accord10      string
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

func (s *Scraper) GetCountry(perfumer string) (string, error) {
	url := fmt.Sprintf("https://www.fragrantica.com/designers/%s.html", perfumer)

	doc, err := s.GetPageBody(url)
	if err != nil {
		return "", fmt.Errorf("Read body failed: %s", err)
	}

	if doc == nil {
		return "", fmt.Errorf("Empty response body")
	}

	var results []string

	doc.Find("main div.col-span-8.col-start-5.md\\:col-span-full").Find("a").Each(func(i int, item *goquery.Selection) {
		results = append(results, strings.TrimSpace(item.Text()))
	})

	if len(results) == 0 {
		return "", fmt.Errorf("No perfumer details found for: %s", perfumer)
	}

	return results[0], nil
}

// Get details of newly found fragrance by parsing html page
// TODO: explore the option of using ML to get the details from the cards...
func (s *Scraper) ParsePageParams(url string) (FragranceParams, error) {
	params := FragranceParams{}

	doc, err := s.GetPageBody(url)
	if err != nil {
		return params, fmt.Errorf("Read body failed: %s", err)
	}

	if doc == nil {
		return params, fmt.Errorf("Empty response body")
	}

	// Gender - for men / for women / for women and men (or vice versa?)
	gender := getGender(doc)
	params.Gender = gender

	// RatingValue & RatingCount
	rVal, rCount := getRatings(doc)
	if rVal != "" {
		params.RatingValue = rVal
		params.RatingCount = int32(rCount)
	}

	// Year
	year, known := getYear(doc)
	if known {
		params.Year = int32(year)
	}

	// Notes [topNotes, middleNotes, baseNotes]
	notes := getNotes(doc)
	params.TopNotes = notes[0]
	params.MiddleNotes = notes[1]
	params.BaseNotes = notes[2]

	// Parfumer1-2 (["unknown"] is received if there were none found)
	perfumers := getPerfumers(doc)
	params.Perfumer1 = perfumers[0]
	if len(perfumers) > 1 {
		params.Perfumer2 = perfumers[1]
	}

	// unpack Accords1-10
	accords := getAccords(doc)
	for i, v := range accords {
		if i >= 10 {
			break
		}

		switch i {
		case 0:
			params.Accord1 = v
		case 1:
			params.Accord2 = v
		case 2:
			params.Accord3 = v
		case 3:
			params.Accord4 = v
		case 4:
			params.Accord5 = v
		case 5:
			params.Accord6 = v
		case 6:
			params.Accord7 = v
		case 7:
			params.Accord8 = v
		case 8:
			params.Accord9 = v
		case 9:
			params.Accord10 = v
		}
	}

	return params, nil
}

func getGender(doc *goquery.Document) string {
	sex := strings.ToLower(strings.TrimSpace(doc.Find("#toptop").Find("span").Text()))

	var gender string

	switch sex {
	case "for men":
		gender = "men"
	case "for women":
		gender = "women"
	default:
		gender = "unisex"
	}

	return gender
}

func getRatings(doc *goquery.Document) (string, int) {
	rating := doc.Find(`[itemprop="aggregateRating"]`)

	value := rating.Find(`[itemprop="ratingValue"]`).Text()

	countSel := rating.Find(`[itemprop="ratingCount"]`)
	countStr, exists := countSel.Attr("content")
	if !exists {
		return "", 0
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", 0
	}

	return value, count
}

func getYear(doc *goquery.Document) (int, bool) {
	s := doc.Find("head > title").Text()

	parts := strings.Fields(s)
	if len(parts) == 0 {
		return 0, false
	}

	last := parts[len(parts)-1]

	year, err := strconv.Atoi(last)
	if err != nil {
		return 0, false
	}

	//basic sanity check
	if year < 1000 || year > 9999 {
		return 0, false
	}

	return year, true
}

func getNotes(doc *goquery.Document) []string {
	var results []string

	// TopNotes
	// #pyramid > div.mx-auto.max-w-md
	var topNotes []string
	doc.Find("#pyramid div.mx-auto.max-w-md").Find("div.flex.flex-wrap.justify-center.items-end.py-3.px-2.pyramid-level-container").Find("span").Each(func(i int, item *goquery.Selection) {
		topNotes = append(topNotes, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	// MiddleNotes
	// #pyramid > div.mx-auto.max-w-xl
	var middleNotes []string
	doc.Find("#pyramid div.mx-auto.max-w-xl").Find("div.flex.flex-wrap.justify-center.items-end.py-3.px-2.pyramid-level-container").Find("span").Each(func(i int, item *goquery.Selection) {
		middleNotes = append(middleNotes, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	// BaseNotes
	// #pyramid > div.mx-auto.max-w-2xl
	var baseNotes []string
	doc.Find("#pyramid div.mx-auto.max-w-2xl").Find("div.flex.flex-wrap.justify-center.items-end.py-3.px-2.pyramid-level-container").Find("span").Each(func(i int, item *goquery.Selection) {
		baseNotes = append(baseNotes, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	results = append(results, strings.Join(topNotes, ", "))
	results = append(results, strings.Join(middleNotes, ", "))
	results = append(results, strings.Join(baseNotes, ", "))

	return results
}

func getPerfumers(doc *goquery.Document) []string {
	var results []string

	// Perfumer1 (optional)
	// Perfumer2 (optional)
	// website has space for 1-4 perfumers, only need 2, set Perfumer1 to unknown if none found
	doc.Find("main div.grid.grid-cols-2.md\\:grid-cols-3.lg\\:grid-cols-4.gap-3.md\\:gap-4").Find("span").Each(func(i int, item *goquery.Selection) {
		results = append(results, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	if len(results) == 0 {
		results = append(results, "unknown")
	}

	return results
}

func getAccords(doc *goquery.Document) []string {
	var results []string

	doc.Find("main div.flex.flex-col.w-full.max-w-\\[280px\\].md\\:max-w-\\[320px\\]").Find("span.truncate").Each(func(i int, item *goquery.Selection) {
		results = append(results, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	return results
}

// func - (to be run periodically (once an hour or so) take one existing frag and update the card and public score (anything else?) if changed
