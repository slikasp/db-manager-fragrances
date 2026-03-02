package cards

import (
	"testing"
)

func TestCropFix(t *testing.T) {
	path := "cards/en/p_c_1.jpeg"

	img, err := cropQR(path)
	if err != nil {
		t.Errorf("QR resolution failed: %v", err)
	}

	err = saveImage(img, "cards/qr/test_crop.jpeg")
	if err != nil {
		t.Errorf("Failed to save: %v", err)
	}

	fixed := fixQR(img)

	err = saveImage(fixed, "cards/qr/test_fix.jpeg")
	if err != nil {
		t.Errorf("Failed to save: %v", err)
	}

	// preprocess?
}

func TestResolveQR(t *testing.T) {
	path := "cards/en/p_c_1.jpeg"

	img, err := cropQR(path)
	if err != nil {
		t.Errorf("Could not crop: %v", err)
	}

	fixed := fixQR(img)

	link, err := decodeGozxing(fixed)
	if err != nil {
		t.Errorf("Could not decode: %v", err)
	}

	stripped, err := stripQuery(link)
	if err != nil {
		t.Errorf("Could not strip URL: %v", err)
	}

	expected := "https://www.fragrantica.com/perfume/Azzaro/Orange-Tonic-1.html"
	if stripped != expected {
		t.Errorf("Unexpected result: %s:%s", stripped, expected)
	}

}
