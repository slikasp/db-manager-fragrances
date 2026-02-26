package cards

import (
	"image"
	"image/color"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

func upscale(src image.Image) image.Image {
	b := src.Bounds()

	borderPx := 5
	scale := 5

	// add white border
	padded := image.NewRGBA(image.Rect(0, 0, b.Dx()+2*borderPx, b.Dy()+2*borderPx))
	draw.Draw(padded, padded.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(padded, image.Rect(borderPx, borderPx, borderPx+b.Dx(), borderPx+b.Dy()), src, b.Min, draw.Src)

	// upscale (nearest neighbor)
	sb := padded.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, sb.Dx()*scale, sb.Dy()*scale))
	xdraw.NearestNeighbor.Scale(out, out.Bounds(), padded, sb, draw.Src, nil)
	return out
}
