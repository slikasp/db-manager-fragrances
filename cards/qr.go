package cards

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/liyue201/goqr"
)

func saveQR(cropped image.Image) error {
	out, err := os.Create("cards/qr/temp.jpeg")
	if err != nil {
		return err
	}
	defer out.Close()

	if err := jpeg.Encode(out, cropped, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}
	return nil
}

func savePNG(img image.Image) error {
	f, err := os.Create("cards/qr/temp.jpeg")
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func resolveQr(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	// Crop rectangle in *image coordinates* (x0,y0) -> (x1,y1)
	rect := image.Rect(1045, 10, 1135, 100) // left, top, right, bottom

	sub, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	if !ok {
		return nil, errors.New("image type does not support cropping")
	}

	cropped := sub.SubImage(rect)

	// --- SAVE TEMP FILE ---
	err = saveQR(cropped)
	if err != nil {
		return nil, err
	}
	// -----------------------

	qrs, err := goqr.Recognize(cropped)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0, len(qrs))
	for _, qr := range qrs {
		results = append(results, string(qr.Payload))
	}

	return results, nil
}
