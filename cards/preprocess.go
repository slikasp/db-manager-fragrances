package cards

import (
	"image"
	"image/color"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

// Chain needed functions here to add specific preprocessing
func PreprocessQR(src image.Image) image.Image {
	return qrOtsuBinarize(src)
}

func qrUpscale(src image.Image) image.Image {
	scale := 4

	// upscale (nearest neighbor)
	pb := src.Bounds()
	up := image.NewRGBA(image.Rect(0, 0, pb.Dx()*scale, pb.Dy()*scale))
	xdraw.NearestNeighbor.Scale(up, up.Bounds(), src, pb, draw.Src, nil)

	return up
}

func qrAddQuietZone(src image.Image) image.Image {
	border := 16

	// add quiet zone (white border)
	b := src.Bounds()
	padded := image.NewRGBA(image.Rect(0, 0, b.Dx()+2*border, b.Dy()+2*border))
	draw.Draw(padded, padded.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(
		padded,
		image.Rect(border, border, border+b.Dx(), border+b.Dy()),
		src,
		b.Min,
		draw.Src,
	)

	return padded
}

func qrOtsuBinarize(img image.Image) image.Image {
	b := img.Bounds()
	g := image.NewGray(b)

	// grayscale + histogram
	var hist [256]uint64
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			gy := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			g.SetGray(x, y, color.Gray{Y: gy})
			hist[gy]++
		}
	}

	// Otsu threshold
	total := uint64(b.Dx() * b.Dy())
	var sum uint64
	for i := 0; i < 256; i++ {
		sum += uint64(i) * hist[i]
	}

	var sumB, wB uint64
	var maxVar float64
	thr := uint8(128)

	for t := 0; t < 256; t++ {
		wB += hist[t]
		if wB == 0 {
			continue
		}
		wF := total - wB
		if wF == 0 {
			break
		}
		sumB += uint64(t) * hist[t]

		mB := float64(sumB) / float64(wB)
		mF := float64(sum-sumB) / float64(wF)
		between := float64(wB) * float64(wF) * (mB - mF) * (mB - mF)

		if between > maxVar {
			maxVar = between
			thr = uint8(t)
		}
	}

	// apply threshold -> pure black/white
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if g.GrayAt(x, y).Y > thr {
				g.SetGray(x, y, color.Gray{Y: 255})
			} else {
				g.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	return g
}
