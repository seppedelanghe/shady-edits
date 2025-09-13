package tuning

import (
	"math/rand"
	"shady-edits/pkg/nodes"
)

func randomOptions(nodeOpts []nodes.NodeOptions, r *rand.Rand) []nodes.NodeOptions {
	candidate := make([]nodes.NodeOptions, 0)
	for _, v := range nodeOpts {
		enabled := r.Float32() < 0.5
		params := make([]nodes.NodeParam, 0)

		for _, param := range v.Params {
			pEnabled := r.Float32() < 0.5
			params = append(params, nodes.NodeParam{
				Enabled: pEnabled,
				Name:    param.Name,
				Value:   r.Float32(),
			})
		}

		candidate = append(candidate, nodes.NodeOptions{
			Name:    v.Name,
			Enabled: enabled,
			Params:  params,
		})
	}
	return candidate
}
