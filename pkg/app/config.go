package app

import (
	"shady-edits/pkg/loss"
	"shady-edits/pkg/nodes"
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

	initialOptions := []nodes.NodeOptions{
		{Name: "Alpha", Enabled: true, Params: []nodes.NodeParam{
			{Name: "Alpha", Enabled: true, Value: 1.0}, // leave fully opaque by default
			{Name: "Blend", Enabled: true, Value: 1.0},
		}},
		{Name: "Exposure", Enabled: true, Params: []nodes.NodeParam{
			{Name: "Exposure", Enabled: true, Value: 0.0}, // neutral
			{Name: "Blend", Enabled: true, Value: 1.0},
		}},
		{Name: "Contrast", Enabled: true, Params: []nodes.NodeParam{
			{Name: "Contrast", Enabled: true, Value: 0.15},
			{Name: "Blend", Enabled: true, Value: 1.0},
		}},
		{Name: "Saturation", Enabled: true, Params: []nodes.NodeParam{
			{Name: "Saturation", Enabled: true, Value: 0.12},
			{Name: "Blend", Enabled: true, Value: 1.0},
		}},
		{Name: "Temperature", Enabled: true, Params: []nodes.NodeParam{
			{Name: "Temperature", Enabled: true, Value: 0.08}, // slight warm
			{Name: "Blend", Enabled: true, Value: 1.0},
		}},
	}

	tuner := tuning.NewRandomGeneticEvolve(initialOptions, 50, 200)

	// tuner := tuning.NewRandomSearch(initialOptions, 1000)

	return &Config{original, target, result, loss.L1LossLinearRGB, tuner, w, h}, nil
}
