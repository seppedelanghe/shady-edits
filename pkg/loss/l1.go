package loss

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// builds a 256-entry LUT mapping 0..255 (sRGB) -> linear value in [0,1].
func buildSRGBToLinearLUT256() [256]float64 {
	var lut [256]float64
	for i := range 256 {
		v := float64(i) / 255.0
		if v <= 0.04045 {
			lut[i] = v / 12.92
		} else {
			lut[i] = math.Pow((v+0.055)/1.055, 2.4)
		}
	}
	return lut
}

// L1LossLinearRGBFast computes mean absolute error in linear RGB space (ignores alpha).
func L1LossLinearRGB(actual, pred *ebiten.Image) float64 {
	b := actual.Bounds()
	if !b.Eq(pred.Bounds()) {
		panic("images must be same size")
	}
	w, h := b.Dx(), b.Dy()
	if w == 0 || h == 0 {
		return 0.0
	}

	lut := buildSRGBToLinearLUT256()

	var sum float64
	inv65535To8 := uint32(8) // shift right by 8 to get top 8 bits (0..255)

	for y := range h {
		py := b.Min.Y + y
		for x := range w {
			px := b.Min.X + x
			r1, g1, b1, _ := actual.At(px, py).RGBA()
			r2, g2, b2, _ := pred.At(px, py).RGBA()

			lr1 := lut[uint8(r1>>inv65535To8)]
			lg1 := lut[uint8(g1>>inv65535To8)]
			lb1 := lut[uint8(b1>>inv65535To8)]

			lr2 := lut[uint8(r2>>inv65535To8)]
			lg2 := lut[uint8(g2>>inv65535To8)]
			lb2 := lut[uint8(b2>>inv65535To8)]

			sum += math.Abs(lr1-lr2) + math.Abs(lg1-lg2) + math.Abs(lb1-lb2)
		}
	}

	// average over pixels and 3 channels
	return sum / (float64(w*h) * 3.0)
}
