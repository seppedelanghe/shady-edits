package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DebugApp struct {
	Pipeline *Pipeline
	Config   *Config

	suspended bool
}

func NewDebugApp(config *Config) DebugApp {
	pipeline := NewDefaultPipeline()
	ebiten.SetWindowSize(config.W*3, config.H)
	ebiten.SetTPS(120)

	return DebugApp{pipeline, config, false}
}

func (dw *DebugApp) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if dw.suspended {
		return nil
	}

	candidate := dw.Config.Tuner.Candidate()
	dw.Config.Result = dw.Pipeline.Run(dw.Config.Original, candidate)
	loss := dw.Config.Loss(dw.Config.Target, dw.Config.Result)
	// fmt.Printf("Loss: %.4f\n", loss)

	if dw.Config.Tuner.Update(loss) {
		nodeOptions := dw.Config.Tuner.NodeOptions()
		dw.Config.Result = dw.Pipeline.Run(dw.Config.Original, nodeOptions)
		loss := dw.Config.Loss(dw.Config.Target, dw.Config.Result)
		fmt.Printf("Found best parameters:\n\t%v\n\tLoss: %.4f", nodeOptions, loss)
		dw.suspended = true
	}

	return nil
}

func (dw *DebugApp) Draw(screen *ebiten.Image) {
	screen.DrawImage(dw.Config.Original, nil)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(dw.Config.W), 0)
	screen.DrawImage(dw.Config.Result, &opts)

	opts2 := ebiten.DrawImageOptions{}
	opts2.GeoM.Translate(float64(dw.Config.W*2), 0)
	screen.DrawImage(dw.Config.Target, &opts2)

	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f, TPS: %.2f", fps, tps))
}

func (dw *DebugApp) Layout(w, h int) (int, int) {
	return w, h
}
func (dw *DebugApp) Run() error {
	return ebiten.RunGame(dw)
}
