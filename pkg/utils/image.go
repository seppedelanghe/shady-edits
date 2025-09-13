package utils

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func CopyImage(src *ebiten.Image) *ebiten.Image {
	w := src.Bounds().Dx()
	h := src.Bounds().Dy()
	dst := ebiten.NewImage(w, h)

	opts := &ebiten.DrawImageOptions{}
	dst.DrawImage(src, opts)

	return dst
}

func CopyRegion(src *ebiten.Image, rect image.Rectangle) *ebiten.Image {
	w, h := rect.Dx(), rect.Dy()
	out := ebiten.NewImage(w, h)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-float64(rect.Min.X), -float64(rect.Min.Y))
	out.DrawImage(src, opts)
	return out
}
