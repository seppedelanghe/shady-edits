package tuning

import (
	"math"
	"math/rand"
	"shady-edits/pkg/nodes"
)

type RandomSearch struct {
	nodeOpts      []nodes.NodeOptions
	candidate     []nodes.NodeOptions
	bestLoss      float64
	iteration     int
	maxIterations int

	r *rand.Rand
}

func NewRandomSearch(initialOpts []nodes.NodeOptions, maxIterations int) *RandomSearch {
	r := rand.New(rand.NewSource(420))
	return &RandomSearch{initialOpts, nil, math.Inf(1), 0, maxIterations, r}
}

func (rs *RandomSearch) Update(loss float64) bool {
	if loss < rs.bestLoss {
		rs.bestLoss = loss
		for i, nodeOpts := range rs.candidate {
			rs.nodeOpts[i] = nodes.NodeOptions{
				Name:    nodeOpts.Name,
				Enabled: nodeOpts.Enabled,
				Params:  nodeOpts.Params,
			}
		}
	}

	rs.iteration++
	rs.candidate = make([]nodes.NodeOptions, 0)
	for _, v := range rs.nodeOpts {
		enabled := rs.r.Float32() < 0.5
		params := make([]nodes.NodeParam, 0)

		for _, param := range v.Params {
			pEnabled := rs.r.Float32() < 0.5
			params = append(params, nodes.NodeParam{
				Enabled: pEnabled,
				Name:    param.Name,
				Value:   rs.r.Float32(),
			})
		}

		rs.candidate = append(rs.candidate, nodes.NodeOptions{
			Name:    v.Name,
			Enabled: enabled,
			Params:  params,
		})
	}

	return rs.iteration >= rs.maxIterations
}

func (rs *RandomSearch) NodeOptions() []nodes.NodeOptions {
	return rs.nodeOpts
}

func (rs *RandomSearch) Candidate() []nodes.NodeOptions {
	return rs.candidate
}
