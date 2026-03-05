package cards

import (
	"testing"
)

func TestCropFix(t *testing.T) {
	path := "cards/en/p_c_24894.jpeg"
	// path := "cards/en/p_c_83846.jpeg"
	// path := "cards/en/p_c_107552.jpeg"

	img, err := cropQR(path)
	if err != nil {
		t.Errorf("QR resolution failed: %v", err)
	}

	err = saveImage(img, "cards/qr/test_crop.jpeg")
	if err != nil {
		t.Errorf("Failed to save: %v", err)
	}

	img = PreprocessQR(img)

	fixed := fixQR(img)

	err = saveImage(fixed, "cards/qr/test_fix.jpeg")
	if err != nil {
		t.Errorf("Failed to save: %v", err)
	}

	// preprocess?
}

func TestResolveQR(t *testing.T) {
	path := "cards/en/p_c_24894.jpeg"
	// path := "cards/en/p_c_83846.jpeg"
	// path := "cards/en/p_c_107552.jpeg"

	img, err := cropQR(path)
	if err != nil {
		t.Errorf("Could not crop: %v", err)
	}

	img = PreprocessQR(img)

	fixed := fixQR(img)

	link, err := decodeGozxing(fixed)
	if err != nil {
		t.Errorf("Could not decode: %v", err)
	}

	stripped, err := stripQuery(link)
	if err != nil {
		t.Errorf("Could not strip URL: %v", err)
	}

	expected := "https://www.fragrantica.com/perfume/Puma/Puma-Yellow-Brasil-Edition-24894.html"
	if stripped != expected {
		t.Errorf("Unexpected result: %s:%s", stripped, expected)
	}

}
