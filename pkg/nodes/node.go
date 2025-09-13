package nodes

import "github.com/hajimehoshi/ebiten/v2"

type Node interface {
	Run(src, dst *ebiten.Image, opts map[string]float32)
}
