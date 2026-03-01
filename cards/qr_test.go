package cards

import (
	"testing"
)

func TestResolveQR(t *testing.T) {
	path := "cards/en/p_c_1.jpeg"

	img, err := cropQR(path)
	if err != nil {
		t.Errorf("QR resolution failed: %v", err)
	}

	link, err := decodeGozxing(img)
	if err != nil {
		t.Errorf("Could not decode: %v", err)
	}

	// if len(link) == 0 {
	// 	t.Error(link)
	// }

	if link != "" {
		t.Error(link)
	}
}

func TestResolveUpscaledQR(t *testing.T) {
	path := "cards/qr/temp.jpeg"

	img, err := decodeImage(path)
	if err != nil {
		t.Errorf("Failed to decode: %v", err)
	}

	pre := preprocessForQR(img)

	err = saveImage(pre, "cards/qr/tempHQ.jpeg")
	if err != nil {
		t.Errorf("Failed to save: %v", err)
	}

	link, err := decodeGozxing(pre)
	if err != nil {
		t.Errorf("Could not decode: %v", err)
	}

	// if len(link) == 0 {
	// 	t.Error(link)
	// }

	if link != "" {
		t.Error(link)
	}
}
