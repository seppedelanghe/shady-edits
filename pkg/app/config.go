package app

import (
	"shady-edits/pkg/loss"
	"shady-edits/pkg/tuning"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Config struct {
	Original *ebiten.Image
	Target   *ebiten.Image
	Result   *ebiten.Image

	Loss  loss.LossFn
	Tuner tuning.ParamTuner

	W, H int
}

func NewConfigFromPaths(originalPath, targetPath string) (*Config, error) {
	original, _, err := ebitenutil.NewImageFromFile(originalPath)
	if err != nil {
		return nil, err
	}

	target, _, err := ebitenutil.NewImageFromFile(targetPath)
	if err != nil {
		return nil, err
	}

	w := target.Bounds().Dx()
	h := target.Bounds().Dy()
	result := ebiten.NewImage(w, h)

	tuner := tuning.NewCoordinateSearch(0.0215, map[string]float32{
		"Alpha":         0.5,
		"AlphaBlend":    1,
		"Contrast":      1,
		"ContrastBlend": 1,
	})

	return &Config{original, target, result, loss.L1Loss, tuner, w, h}, nil
}
