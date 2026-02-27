package cards

import (
	"errors"
	"image"
	"image/jpeg"
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

	err = saveQR(cropped)
	if err != nil {
		return nil, err
	}

	return cropped, nil
}

func decodeQR(qr image.Image) ([]string, error) {
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
