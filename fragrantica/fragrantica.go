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

	// ID - from DB
	// URL - from DB

	// Name - from URL
	// Brand - from URL
	// FragranticaID - from URL

	// Country - open brand page and look there
	// #app > main > div > div.col-span-12.sm\:col-span-9.lg\:col-span-9.lg\:pl-1.lg\:mt-1 > div.bg-white.dark\:bg-zinc-900.dark\:text-zinc-100.p-2.md\:p-4.pb-20.rounded-md.relative > div.grid.grid-cols-2.gap-4.md\:gap-6.md\:mb-2 > div.text-center > p > a
	// <a itemprop="url" href="https://www.fragrantica.com/designers/Amouage.html" class="inline-block hover:opacity-80 transition-opacity"><span class="block text-xs md:text-base text-zinc-800 dark:text-zinc-100 font-medium mb-2" itemprop="name"> Amouage </span><span class="inline-block bg-white p-1 rounded-lg"><img class="max-w-14 md:max-w-[150px] block" itemprop="logo" src="https://fimgs.net/mdimg/dizajneri/m.122.jpg" alt="Amouage logo"></span></a>
	// then
	// #app > main > div > div.col-span-12.sm\:col-span-9.lg\:col-span-9.lg\:pl-1.lg\:mt-1 > div.grid.grid-cols-1.gap-4.bg-white.dark\:bg-zinc-800.dark\:text-zinc-100.p-2.md\:p-4.rounded-md > div.grid.grid-cols-1.md\:grid-cols-3.gap-4.md\:gap-6 > div.md\:col-span-1 > div > div.col-span-8.col-start-5.md\:col-span-full
	// <a href="/country/Oman.html" class="font-bold text-teal-900 dark:text-teal-500">Oman</a>

	// Gender - for men / for women / for men and women
	// #toptop
	// <h1 itemprop="name" class="text-2xl md:text-[2.5rem] font-light tracking-tight text-center md:text-left text-zinc-800 dark:text-zinc-100"> Reflection Man Amouage
	// <span class="text-lg md:text-2xl text-blue-600 dark:text-blue-400 whitespace-nowrap">for men</span>
	// </h1>

	// RatingValue -
	// RatingCount -
	// #app > main > div > div.col-span-12.sm\:col-span-9.lg\:col-span-9.lg\:pl-1.lg\:mt-1 > div.bg-white.dark\:bg-zinc-900.dark\:text-zinc-100.p-2.md\:p-4.pb-20.rounded-md.relative > div.mb-8 > div.mt-4.flex.flex-wrap.items-center.justify-center.gap-2
	// <p class="text-xs sm:text-sm text-zinc-500 dark:text-zinc-400" itemprop="aggregateRating" itemtype="http://schema.org/AggregateRating" itemscope=""> Perfume rating&nbsp;
	// <span itemprop="ratingValue" class="font-semibold text-teal-600 dark:text-teal-400">4.40</span>&nbsp;out of&nbsp;
	// <span itemprop="bestRating">5</span>&nbsp;with&nbsp;
	// <span itemprop="ratingCount" content="9730" class="font-semibold">9,730</span>&nbsp;votes </p>

	// Year -
	// head > title

	// TopNotes -
	// #pyramid > div.relative.bg-linear-to-br.from-white.to-zinc-50\/80.dark\:from-zinc-800.dark\:to-zinc-800\/80.rounded-xl.shadow-sm.shadow-zinc-300\/50.dark\:shadow-black\/20.overflow-hidden > div.p-5 > div > div.mt-6.space-y-1 > div.mx-auto.max-w-md
	// <a href="https://www.fragrantica.com/notes/Rosemary-49.html" class="group relative flex flex-col items-center text-center pyramid-note-link" style="opacity: 0.915348;"><div class="relative"><img loading="lazy" src="https://fimgs.net/mdimg/sastojci/t.49.jpg" class="rounded-md shadow-xs ring-1 ring-zinc-200/20 dark:ring-zinc-700/30 transition-all duration-300 ease-out group-hover:scale-110 group-hover:shadow-md group-hover:ring-teal-400/40 dark:group-hover:ring-teal-500/40" alt="Rosemary" style="width: 3.4rem;"></div><span class="pyramid-note-label mt-1.5 text-[11px] sm:text-sm font-medium text-zinc-600 dark:text-zinc-200 group-hover:text-teal-600 dark:group-hover:text-teal-400 transition-colors duration-200 whitespace-nowrap"> Rosemary </span></a>	// MiddleNotes -
	// <a href="https://www.fragrantica.com/notes/Pink-Pepper-91.html" class="group relative flex flex-col items-center text-center pyramid-note-link" style="opacity: 0.765573;"><div class="relative"><img loading="lazy" src="https://fimgs.net/mdimg/sastojci/t.91.jpg" class="rounded-md shadow-xs ring-1 ring-zinc-200/20 dark:ring-zinc-700/30 transition-all duration-300 ease-out group-hover:scale-110 group-hover:shadow-md group-hover:ring-teal-400/40 dark:group-hover:ring-teal-500/40" alt="Pink Pepper" style="width: 2.625rem;"></div><span class="pyramid-note-label mt-1.5 text-[11px] sm:text-sm font-medium text-zinc-600 dark:text-zinc-200 group-hover:text-teal-600 dark:group-hover:text-teal-400 transition-colors duration-200 whitespace-nowrap"> Pink Pepper </span></a>
	// <a href="https://www.fragrantica.com/notes/Petitgrain-3.html" class="group relative flex flex-col items-center text-center pyramid-note-link" style="opacity: 0.747551;"><div class="relative"><img loading="lazy" src="https://fimgs.net/mdimg/sastojci/t.3.jpg" class="rounded-md shadow-xs ring-1 ring-zinc-200/20 dark:ring-zinc-700/30 transition-all duration-300 ease-out group-hover:scale-110 group-hover:shadow-md group-hover:ring-teal-400/40 dark:group-hover:ring-teal-500/40" alt="Petitgrain" style="width: 2.55rem;"></div><span class="pyramid-note-label mt-1.5 text-[11px] sm:text-sm font-medium text-zinc-600 dark:text-zinc-200 group-hover:text-teal-600 dark:group-hover:text-teal-400 transition-colors duration-200 whitespace-nowrap"> Petitgrain </span></a>
	// #pyramid > div.relative.bg-linear-to-br.from-white.to-zinc-50\/80.dark\:from-zinc-800.dark\:to-zinc-800\/80.rounded-xl.shadow-sm.shadow-zinc-300\/50.dark\:shadow-black\/20.overflow-hidden > div.p-5 > div > div.mt-6.space-y-1 > div.mx-auto.max-w-xl
	// etc
	// BaseNotes -
	// #pyramid > div.relative.bg-linear-to-br.from-white.to-zinc-50\/80.dark\:from-zinc-800.dark\:to-zinc-800\/80.rounded-xl.shadow-sm.shadow-zinc-300\/50.dark\:shadow-black\/20.overflow-hidden > div.p-5 > div > div.mt-6.space-y-1 > div.mx-auto.max-w-2xl
	// etc

	// Parfumer1 -
	// Parfumer2 -
	// #app > main > div > div.col-span-12.sm\:col-span-9.lg\:col-span-9.lg\:pl-1.lg\:mt-1 > div.bg-white.dark\:bg-zinc-900.dark\:text-zinc-100.p-2.md\:p-4.pb-20.rounded-md.relative > div:nth-child(9) > div.grid.grid-cols-2.md\:grid-cols-3.lg\:grid-cols-4.gap-3.md\:gap-4
	// <a href="/noses/Lucas_Sieuzac.html" class="group flex items-center gap-2 md:gap-3 px-3 md:px-4 py-2 md:py-3 bg-linear-to-br from-white to-zinc-50 dark:from-zinc-800 dark:to-zinc-800/80 rounded-xl border border-zinc-200/60 dark:border-zinc-700/50 shadow-sm hover:shadow-md dark:shadow-black/20 hover:border-teal-300 dark:hover:border-teal-600/50 transition-shadow duration-150"><div class="relative shrink-0"><img src="https://frgs.me/mdimg/nosevi/fit.168.jpg" class="w-10 h-10 md:w-12 md:h-12 rounded-full object-cover ring-2 ring-zinc-200 dark:ring-zinc-600 group-hover:ring-teal-400 dark:group-hover:ring-teal-500 transition-shadow duration-150" alt="Lucas Sieuzac"></div><span class="text-sm font-medium text-zinc-700 dark:text-zinc-200 group-hover:text-teal-600 dark:group-hover:text-teal-400 transition-colors"> Lucas Sieuzac </span></a></div>
	// can be 1-4 perfumers, only need 2, handle nil

	// Accord1
	// Accord2
	// Accord3
	// Accord4
	// Accord5
	// accords := getAccords(doc)

}

func getAccords(doc *goquery.Document) []string {
	var results []string

	doc.Find("main div.flex.flex-col.w-full.max-w-\\[280px\\].md\\:max-w-\\[320px\\]").Each(func(i int, s *goquery.Selection) {
		s.Find("span.truncate").Each(func(j int, item *goquery.Selection) {
			results = append(results, item.Text())
		})
	})

	return results
}

// func - (to be run periodically (once an hour or so) take one existing frag and update the card and public score (anything else?) if changed
