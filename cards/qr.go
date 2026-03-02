package cards

import (
	"errors"
	"image"
	"image/jpeg"
	"net/url"
	"os"

	"github.com/liyue201/goqr"

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
		return nil, err
	}

	// Crop rectangle in *image coordinates* (x0,y0) -> (x1,y1)
	rect := image.Rect(1045, 10, 1135, 100) // left, top, right, bottom

	sub, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	if !ok {
		return nil, errors.New("Image type does not support cropping")
	}

	cropped := sub.SubImage(rect)

	return cropped, nil
}

// Trim duplicated rows/columns from the QR code
// TODO: add preprocessing so we only have black and white pixels only
// TODO: update the code so it detects duplicate rows/columns
func fixQR(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	drop := []int{1, 4, 8, 9, 12, 15, 19, 20, 23, 26, 31, 32, 34, 38, 39, 42, 45, 49, 50, 53, 56, 60, 61, 64, 68, 69, 72, 75, 79, 80, 83, 87, 88}

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

// Could not decode without upscaling the image
func decodeGoqr(qr image.Image) ([]string, error) {
	qrs, err := goqr.Recognize(qr)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0, len(qrs))
	for _, qr := range qrs {
		results = append(results, string(qr.Payload))
	}

	return results, nil
}

// Decode the QR from an umage using gozxing
func decodeGozxing(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	reader := qrcode.NewQRCodeReader()

	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}

	result, err := reader.Decode(bmp, hints)
	if err != nil {
		return "", err
	}

	return result.GetText(), nil
}

// Drop parameters from the links:
//
// https://www.fragrantica.com/perfume/Azzaro/Orange-Tonic-1.html?utm_source=qr-code&utm_medium=social-card
//
//	->
//
// https://www.fragrantica.com/perfume/Azzaro/Orange-Tonic-1.html
func stripQuery(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	u.RawQuery = ""
	u.Fragment = ""

	return u.String(), nil
}
