package app

import (
	"shady-edits/pkg/nodes"
	"shady-edits/pkg/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

var alphaShader = `
//kage:unit pixels

package main

var Alpha float
var AlphaBlend float

func Fragment(dstPos vec4, srcPos vec2) vec4 {
	color := imageSrc0UnsafeAt(srcPos)
	return vec4(color.rgb * Alpha, AlphaBlend)
}
`

var contrastShader = `
//kage:unit pixels

package main

var Contrast float
var ContrastBlend float

func Fragment(dstPos vec4, srcPos vec2) vec4 {
	color := imageSrc0UnsafeAt(srcPos)
	color.rgb = (color.rgb - 0.5) * Contrast + 0.5;
	return vec4(color.rgb, ContrastBlend)
}

`

type Pipeline struct {
	nodes []nodes.Node

	src, dst *ebiten.Image
}

func NewDefaultPipeline() *Pipeline {
	nodes := []nodes.Node{
		nodes.NewShaderNode([]byte(alphaShader)),
		nodes.NewShaderNode([]byte(contrastShader)),
	}
	return &Pipeline{nodes, nil, nil}
}

func (p *Pipeline) Run(image *ebiten.Image, opts map[string]float32) *ebiten.Image {
	p.src = utils.CopyImage(image)
	if p.dst == nil {
		p.dst = ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
	} else {
		p.dst.Clear()
	}

	var src *ebiten.Image
	for _, node := range p.nodes {
		node.Run(p.src, p.dst, opts)
		src = p.dst
		p.dst = p.src
		p.src = src
	}

	if len(p.nodes)%2 == 0 {
		return p.dst
	}
	return p.src
}
