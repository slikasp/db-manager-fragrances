package cards

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/url"
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func saveImage(cropped image.Image, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := jpeg.Encode(out, cropped, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}
	return nil
}

func decodeImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Crop the QR code from the card image.
// It has to be in the pixel perfect position, from my checks it always is.
// Will need to crop few fixels around it if it's not or add detection.
func cropQR(filePath string) (image.Image, error) {
	img, err := decodeImage(filePath)
	if err != nil {
		return nil, fmt.Errorf("Could not decode image %s: %w", filePath, err)
	}

	// Crop rectangle in *image coordinates* (x0,y0) -> (x1,y1)
	rect := image.Rect(1041, 6, 1139, 104) // left, top, right, bottom

	sub, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	if !ok {
		return nil, fmt.Errorf("Image type does not support cropping: %s", filePath)
	}

	cropped := sub.SubImage(rect)

	return cropped, nil
}

// Rreturns indices (0-based) of rows to remove.
// Assumptions:
// - duplicates only occur as adjacent pairs (y and y+1)
// - each duplicated row appears only once
// - black/white, but we just compare exact pixels (RGBA) for safety
func duplicateLines(img image.Image) []int {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	var remove []int
	for y := 0; y < h-1; y++ {
		if linesEqual(img, b.Min.X, b.Min.Y+y, w) {
			// row y and y+1 are identical -> remove the second one
			remove = append(remove, y+1)
			y++ // skip next row since we just matched the pair
		}
	}
	return remove
}

// Compares row at (minX, y0) with the next row (minX, y0+1), width w.
func linesEqual(img image.Image, minX, y0, w int) bool {
	y1 := y0 + 1
	for x := 0; x < w; x++ {
		c0 := img.At(minX+x, y0)
		c1 := img.At(minX+x, y1)
		if c0 != c1 {
			// Note: color.Color is an interface; direct != is usually fine here,
			// but to be absolutely consistent across color models, compare RGBA:
			r0, g0, b0, a0 := c0.RGBA()
			r1, g1, b1, a1 := c1.RGBA()
			if r0 != r1 || g0 != g1 || b0 != b1 || a0 != a1 {
				return false
			}
		}
	}
	return true
}

// Trim duplicated rows/columns from the QR code
func fixQR(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	// Rows and columns are duplicated on the same offset
	drop := duplicateLines(src)

	// Build lookup set
	dropSet := make(map[int]struct{}, len(drop))
	for _, v := range drop {
		dropSet[v] = struct{}{}
	}

	// Calculate new size
	newW, newH := 0, 0
	for x := 0; x < w; x++ {
		if _, remove := dropSet[x]; !remove {
			newW++
		}
	}
	for y := 0; y < h; y++ {
		if _, remove := dropSet[y]; !remove {
			newH++
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// Copy pixels
	dy := 0
	for y := 0; y < h; y++ {
		if _, remove := dropSet[y]; remove {
			continue
		}
		dx := 0
		for x := 0; x < w; x++ {
			if _, remove := dropSet[x]; remove {
				continue
			}
			dst.Set(dx, dy, src.At(b.Min.X+x, b.Min.Y+y))
			dx++
		}
		dy++
	}

	return dst
}

// Decode the QR from an umage using gozxing
func decodeGozxing(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("Failed converting image to bitmap: %w", err)
	}

	reader := qrcode.NewQRCodeReader()

	// hints := map[gozxing.DecodeHintType]interface{}{
	// 	gozxing.DecodeHintType_TRY_HARDER: true,
	// }
	hints := map[gozxing.DecodeHintType]interface{}{}

	result, err := reader.Decode(bmp, hints)
	if err != nil {
		return "", fmt.Errorf("Failed decoding bitmap: %w", err)
	}

	return result.GetText(), nil
}

// Drop parameters from the links:
func stripQuery(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	u.RawQuery = ""
	u.Fragment = ""

	return u.String(), nil
}

// Gets the link to the fragrance from it's card image
func GetLinkFromCard(cardPath string) (string, error) {
	// Crop QR from card image
	img, err := cropQR(cardPath)
	if err != nil {
		return "", fmt.Errorf("Failed to crop image %s: %w", cardPath, err)
	}

	// ONLY FOR TESTING
	// saveImage(img, "cards/cards/qr/temp_crop.jpeg")

	img = PreprocessQR(img)

	// ONLY FOR TESTING
	// saveImage(img, "cards/cards/qr/temp_prep.jpeg")

	// strip duplicate rows/columns from QR image
	img = fixQR(img)

	// ONLY FOR TESTING
	// saveImage(img, "cards/cards/qr/temp_fixed.jpeg")

	// decode QR code
	link, err := decodeGozxing(img)
	if err != nil {
		return "", fmt.Errorf("Failed decoding the QR code from %s: %w", cardPath, err)
	}

	// strip query parameters from link
	link, err = stripQuery(link)
	if err != nil {
		return "", fmt.Errorf("Failed stripping query parameters from URL %s: %w", link, err)
	}

	// normalize string
	link = strings.ToLower(link)

	return link, nil
}
