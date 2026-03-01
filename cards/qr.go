package cards

import (
	"errors"
	"image"
	"image/jpeg"
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

	err = saveImage(cropped, "cards/qr/temp.jpeg")
	if err != nil {
		return nil, err
	}

	return cropped, nil
}

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
