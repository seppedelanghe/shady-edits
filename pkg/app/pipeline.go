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
var Blend float

func Fragment(dstPos vec4, srcPos vec2) vec4 {
	color := imageSrc0UnsafeAt(srcPos)
	return vec4(color.rgb * Alpha, Blend)
}
`

var contrastShader = `
//kage:unit pixels

package main

// -1 to 1
var Contrast float

// 0 to 1
var Blend float

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
	outRgb := mix(src.rgb, sRGB, Blend)

	return vec4(outRgb, src.a)
}
`

var exposureShader = `
//kage:unit pixels
package main

// Exposure in stops. 1.0 == +1 stop (x2), -1.0 == -1 stop (x0.5)
var Exposure float

// 0 to 1
var Blend float
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
	linear := sRGBToLinear(src.rgb)

	gain := pow(2.0, Exposure)
	adjustedLinear := linear * gain
	adjustedLinear = clamp(adjustedLinear, vec3(0.0), vec3(1.0))

	sRGB := linearTosRGB(adjustedLinear)
	outRgb := mix(src.rgb, sRGB, Blend)
	return vec4(outRgb, src.a)
}
`

var saturationShader = `
//kage:unit pixels
package main
// -1 to 1 : -1 = grayscale, 0 = original, 1 = doubled saturation
var Saturation float
// 0 to 1
var Blend float
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
	linear := sRGBToLinear(src.rgb)
	// Rec.709 luma in linear space
	luma := dot(linear, vec3(0.2126, 0.7152, 0.0722))
	gray := vec3(luma, luma, luma)
	// factor: 1 -> original, 0 -> grayscale, >1 -> boost saturation
	factor := 1.0 + Saturation
	adjustedLinear := gray + (linear - gray) * factor
	adjustedLinear = clamp(adjustedLinear, vec3(0.0), vec3(1.0))
	sRGB := linearTosRGB(adjustedLinear)
	outRgb := mix(src.rgb, sRGB, Blend)
	return vec4(outRgb, src.a)
}
`

var temperatureShader = `
//kage:unit pixels
package main

// -1 to 1 : positive warms (more yellow/red), negative cools (more blue)
var Temperature float

// 0 to 1
var Blend float

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

// Approximate color temperature (Kelvin) to linear RGB using a known approximation.
// Input: temperature in Kelvin (e.g., 6500). Output: RGB in 0..1 range.
func kelvinToRGB(k float) vec3 {
	t := k / 100.0
	var r float
	var g float
	var b float
	// Red
	if t <= 66.0 {
		r = 1.0 // 255/255
	} else {
		// 329.698727446 * (t - 60) ^ -0.1332047592
		r = 329.698727446 * pow(t-60.0, -0.1332047592) / 255.0
	}
	// Green
	if t <= 66.0 {
		g = 99.4708025861 * log(t) - 161.1195681661
		g = g / 255.0
	} else {
		g = 288.1221695283 * pow(t-60.0, -0.0755148492) / 255.0
	}
	// Blue
	if t >= 66.0 {
		b = 1.0
	} else if t <= 19.0 {
		b = 0.0
	} else {
		b = 138.5177312231 * log(t-10.0) - 305.0447927307
		b = b / 255.0
	}
	// clamp components
	r = clamp(r, 0.0, 1.0)
	g = clamp(g, 0.0, 1.0)
	b = clamp(b, 0.0, 1.0)
	return vec3(r, g, b)
}

func Fragment(dstPos vec4, srcPos vec2) vec4 {
	src := imageSrc0UnsafeAt(srcPos)
	// Map Temperature (-1..1) to a Kelvin range around neutral 6500K.
	// Positive Temperature -> warmer (lower Kelvin), Negative -> cooler (higher Kelvin).
	kelv := 6500.0 - Temperature * 2000.0
	kelv = clamp(kelv, 1000.0, 40000.0)

	// Compute RGB white points and derive per-channel gains relative to neutral 6500K.
	targetRGB := kelvinToRGB(kelv)
	neutralRGB := kelvinToRGB(6500.0)

	// Avoid division by zero; neutralRGB components expected > 0 but clamp defensively.
	neutralRGB = max(neutralRGB, vec3(0.0001))
	gains := targetRGB / neutralRGB
	linear := sRGBToLinear(src.rgb)
	adjustedLinear := linear * gains
	adjustedLinear = clamp(adjustedLinear, vec3(0.0), vec3(1.0))

	sRGB := linearTosRGB(adjustedLinear)
	outRgb := mix(src.rgb, sRGB, Blend)
	return vec4(outRgb, src.a)
}
`

type Pipeline struct {
	nodes map[string]nodes.Node

	src, dst *ebiten.Image
}

func NewDefaultPipeline() *Pipeline {
	nodes := map[string]nodes.Node{
		"Alpha":       nodes.NewShaderNode([]byte(alphaShader)),
		"Contrast":    nodes.NewShaderNode([]byte(contrastShader)),
		"Saturation":  nodes.NewShaderNode([]byte(saturationShader)),
		"Exposure":    nodes.NewShaderNode([]byte(exposureShader)),
		"Temperature": nodes.NewShaderNode([]byte(temperatureShader)),
	}
	return &Pipeline{nodes, nil, nil}
}

func (p *Pipeline) Run(image *ebiten.Image, opts []nodes.NodeOptions) *ebiten.Image {
	p.src = utils.CopyImage(image)
	if p.dst == nil {
		p.dst = ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
	} else {
		p.dst.Clear()
	}

	var src *ebiten.Image
	for _, nodeOpts := range opts {
		if !nodeOpts.Enabled {
			continue
		}

		node := p.nodes[nodeOpts.Name]
		node.Run(p.src, p.dst, nodeOpts)
		src = p.dst
		p.dst = p.src
		p.src = src
	}

	if len(p.nodes)%2 == 0 {
		return p.dst
	}
	return p.src
}
