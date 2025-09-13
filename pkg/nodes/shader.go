package nodes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ShaderNode struct {
	shader *ebiten.Shader
}

func NewShaderNode(shaderCode []byte) *ShaderNode {
	shader, err := ebiten.NewShader(shaderCode)
	if err != nil {
		panic(err)
	}

	return &ShaderNode{shader}
}

func (n *ShaderNode) Run(src, dst *ebiten.Image, opts map[string]float32) {
	uniforms := make(map[string]any)
	for k, v := range opts {
		uniforms[k] = v
	}

	shaderOpts := ebiten.DrawRectShaderOptions{}
	shaderOpts.Images[0] = src
	shaderOpts.Uniforms = uniforms
	dst.DrawRectShader(dst.Bounds().Dx(), dst.Bounds().Dy(), n.shader, &shaderOpts)
}
