package loss

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func L1Loss(actual, pred *ebiten.Image) float64 {
	w := actual.Bounds().Dx()
	h := actual.Bounds().Dy()
	if w != pred.Bounds().Dx() || h != pred.Bounds().Dy() {
		panic("images must be same size to do prediction")
	}

	values := make([]int16, w*h)
	for y := range h {
		for x := range w {
			r1, g1, b1, a1 := actual.At(x, y).RGBA()
			r2, g2, b2, a2 := pred.At(x, y).RGBA()
			values[y*x+y] = int16((r1 - r2) + (g1 - g2) + (b1 - b2) + (a1 - a2))
		}
	}

	var sum float64
	for _, v := range values {
		sum += math.Abs(float64(v))
	}
	return sum / float64(w*h*4)
}
