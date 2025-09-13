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

// -1 to 1
var Contrast float

// 0 to 1
var ContrastBlend float

func sRGBLinearFloat(v float) float {
	if v <= 0.04045 {
		return v / 12.92
	}
	return pow((v+0.055)/1.055, 2.4)
}

func linearsRGBFloat(v float) float {
	if v <= 0.0031308 {
		return v * 12.92
	}
	return 1.055 * pow(v, 1.0 / 2.4) - 0.055
}

func sRGBToLinear(c vec3) vec3 {
	return vec3(sRGBLinearFloat(c.r), sRGBLinearFloat(c.g), sRGBLinearFloat(c.b))
}

func linearTosRGB(c vec3) vec3 {
	return vec3(linearsRGBFloat(c.r), linearsRGBFloat(c.g), linearsRGBFloat(c.b))
}

func Fragment(dstPos vec4, srcPos vec2) vec4 {
	src := imageSrc0UnsafeAt(srcPos)
	C := 1.0 + Contrast // Normalize

	linear := sRGBToLinear(src.rgb)

	luma := dot(linear, vec3(0.2126, 0.7152, 0.0722)) // Rec.709 luma
	lumaAdj := (luma - 0.5) * C + 0.5
	lumaAdj = clamp(lumaAdj, 0.0, 1.0)

	scale := 0.0
	if luma > 0.0 {
		scale = lumaAdj / luma
	}
	adjusted := src.rgb * scale
	adjusted = clamp(adjusted, vec3(0.0), vec3(1.0))

	sRGB := linearTosRGB(adjusted)
	outRgb := mix(src.rgb, sRGB, ContrastBlend)

	return vec4(outRgb, src.a)
}

`

type Pipeline struct {
	nodes []nodes.Node

	src, dst *ebiten.Image
}

func NewDefaultPipeline() *Pipeline {
	nodes := []nodes.Node{
		// nodes.NewShaderNode([]byte(alphaShader)),
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
