package fragrantica

import (
	"slices"
	"testing"
)

func TestReadAndParse(t *testing.T) {
	url := "https://www.fragrantica.com/perfume/Amouage/Reflection-Man-920.html"

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

	accords := getAccords(doc)
	expectedAccords := []string{"woody", "white floral", "aromatic", "powdery", "fresh spicy"}
	if !slices.Equal(accords[:5], expectedAccords) {
		t.Errorf("got %v, want %v", accords, expectedAccords)
	}

}
