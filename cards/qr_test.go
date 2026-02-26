package cards

import (
	"testing"
)

func TestResolveQR(t *testing.T) {
	path := "cards/en/p_c_1.jpeg"

	links, err := resolveQr(path)
	if err != nil {
		t.Errorf("QR resolution failed: %v", err)
	}

	if len(links) == 0 {
		t.Error(links)
	}
}
