package loss

import "github.com/hajimehoshi/ebiten/v2"

type LossFn func(actual, pred *ebiten.Image) float64
