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

func (n *ShaderNode) Run(src, dst *ebiten.Image, opts NodeOptions) {
	uniforms := make(map[string]any)
	for _, param := range opts.Params {
		if param.Enabled {
			uniforms[param.Name] = param.Value
		}
	}

	shaderOpts := ebiten.DrawRectShaderOptions{}
	shaderOpts.Images[0] = src
	shaderOpts.Uniforms = uniforms
	dst.DrawRectShader(dst.Bounds().Dx(), dst.Bounds().Dy(), n.shader, &shaderOpts)
}
