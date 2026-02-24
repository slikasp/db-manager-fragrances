package cards

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	log := Log{}

	log.card(1, true)

	if len(log.Found) != 1 {
		t.Error("'Found' not logged")
	}

	log.card(2, false)

	if len(log.NotFound) != 1 {
		t.Error("'NotFound' not logged")
	}

	log.card(3, true)

	if log.Checked != 3 {
		t.Error("'Checked' wrong number")
	}
}

func TestMakers(t *testing.T) {
	expectedUrl := "https://fimgs.net/mdimg/perfume-social-cards/en-p_c_1.jpeg"
	url := makeCardURL(1)

	if url != expectedUrl {
		t.Errorf("Bad URL: %s:%s", url, expectedUrl)
	}

	expectedPath := "cards/en/p_c_1.jpeg"
	path := makeFilePath(1)

	if path != expectedPath {
		t.Errorf("Bad URL: %s:%s", path, expectedPath)
	}
}

func TestDownload(t *testing.T) {
	card1, err := downloadCard(1)
	if err != nil {
		t.Errorf("Download failed: %v", err)
	}
	_, err = os.Stat(card1.Image)
	if err != nil {
		t.Errorf("Image not there: %v", err)
	}

	card2, err := downloadCard(2)
	if err == nil {
		t.Error("Download should have failed")
	}

	if card2.HasCard {
		t.Error("Card should not exist.")
	}
}
