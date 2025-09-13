package utils

import (
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
