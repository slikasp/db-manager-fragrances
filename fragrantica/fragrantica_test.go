package fragrantica

import (
	"slices"
	"testing"
)

func TestReadAndParse(t *testing.T) {
	url := "https://www.fragrantica.com/perfume/Amouage/Reflection-Man-920.html"
	// url := "https://www.fragrantica.com/perfume/jean-paul-gaultier/le-male-pride-2024-90393.html"
	// url := "https://www.fragrantica.com/perfume/guerlain/neroli-outrenoir-2024-89177.html"

	scraper, err := NewScraper()
	if err != nil {
		t.Fatalf("Failed creating scraper: %s", err)
	}

	doc, err := scraper.GetPageBody(url)
	if err != nil {
		t.Errorf("Read body failed: %s", err)
	}

	if doc == nil {
		t.Errorf("Empty response body")
	}

	// use ParseParams when finished and compare to a pre-made struct

	accords := getAccords(doc)
	expectedAccords := []string{"woody", "white floral", "aromatic", "powdery", "fresh spicy"}
	if !slices.Equal(accords[:5], expectedAccords) {
		t.Errorf("got %v, want %v", accords, expectedAccords)
	}

	perfumers := getPerfumers(doc)
	expectedPerfumers := []string{"lucas sieuzac"}
	if len(perfumers) != 1 {
		t.Errorf("Expected %d perfumers, got %d", 1, len(perfumers))
	}
	if !slices.Equal(perfumers, expectedPerfumers) {
		t.Errorf("got %v, want %v", perfumers, expectedPerfumers)
	}

}
