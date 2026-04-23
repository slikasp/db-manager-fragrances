package fragrantica

import (
	"testing"
)

func TestReadAndParse(t *testing.T) {

	// url := "https://www.fragrantica.com/perfume/jean-paul-gaultier/le-male-pride-2024-90393.html"
	// url := "https://www.fragrantica.com/perfume/guerlain/neroli-outrenoir-2024-89177.html"

	scraper, err := NewScraper()
	if err != nil {
		t.Fatalf("Failed creating scraper: %s", err)
	}

	reflectionUrl := "https://www.fragrantica.com/perfume/Amouage/Reflection-Man-920.html"

	doc, err := scraper.GetPageBody(reflectionUrl)
	if err != nil {
		t.Errorf("Read body failed: %s", err)
	}

	if doc == nil {
		t.Errorf("Empty response body")
	}

	params, err := ParsePageParams(reflectionUrl)
	if err != nil {
		t.Errorf("Page parsing failed: %s", err)
	}

	reflectionParams := FragranceParams{
		Gender:      "men",
		RatingValue: "4.39",
		RatingCount: 9745,
		Year:        2006,
		TopNotes:    "rosemary, pink pepper, petitgrain",
		MiddleNotes: "jasmine, neroli, orris root, ylang-ylang",
		BaseNotes:   "sandalwood, cedar, vetiver, patchouli",
		Perfumer1:   "lucas sieuzac",
		Accord1:     "woody",
		Accord2:     "white floral",
		Accord3:     "aromatic",
		Accord4:     "powdery",
		Accord5:     "fresh spicy",
		Accord6:     "iris",
		Accord7:     "citrus",
		Accord8:     "earthy",
		Accord9:     "floral",
		Accord10:    "warm spicy",
	}

	if params != reflectionParams {
		t.Errorf("got %v\nwant %v", params, reflectionParams)
	}

	demeterUrl := "https://www.fragrantica.com/perfume/demeter-fragrance/jelly-belly-wild-blackberry-peach-cobbler-7702.html"

	doc, err = scraper.GetPageBody(demeterUrl)
	if err != nil {
		t.Errorf("Read body failed: %s", err)
	}

	if doc == nil {
		t.Errorf("Empty response body")
	}

	params, err = ParsePageParams(demeterUrl)
	if err != nil {
		t.Errorf("Page parsing failed: %s", err)
	}

	demeterParams := FragranceParams{
		Gender:      "women",
		RatingValue: "3.85",
		RatingCount: 26,
		TopNotes:    "amalfi lemon",
		MiddleNotes: "blackberry",
		BaseNotes:   "peach",
		Perfumer1:   "unknown",
		Accord1:     "fruity",
		Accord2:     "citrus",
		Accord3:     "sweet",
		Accord4:     "aromatic",
		Accord5:     "violet",
	}

	if params != demeterParams {
		t.Errorf("got %v\nwant %v", params, demeterParams)
	}
}

func TestGetCountry(t *testing.T) {
	amouage := "amouage"

	expectedCountry := "Oman"

	country, err := GetCountry(amouage)
	if err != nil {
		t.Errorf("Failed getting country: %s", err)
	}

	if country != expectedCountry {
		t.Errorf("got %v\nwant %v", country, expectedCountry)
	}
}
