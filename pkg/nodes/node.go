package nodes

import "github.com/hajimehoshi/ebiten/v2"

type NodeParam struct {
	Enabled bool
	Name    string
	Value   float32
}

type NodeOptions struct {
	Enabled bool
	Name    string
	Params  []NodeParam
}

type Node interface {
	Run(src, dst *ebiten.Image, opts NodeOptions)
}
