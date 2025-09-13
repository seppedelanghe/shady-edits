package tuning

import (
	"math"
	"math/rand"
	"shady-edits/pkg/nodes"
	"shady-edits/pkg/utils"
)

type RandomGeneticEvolve struct {
	nodeOpts []nodes.NodeOptions

	candidates     [][]nodes.NodeOptions
	candidatesLoss []float64
	iteration      int
	generation     int

	populations int
	generations int

	r  *rand.Rand
	pb utils.ProgressBar
}

func NewRandomGeneticEvolve(initialParams []nodes.NodeOptions, populations, generations int) *RandomGeneticEvolve {
	candidates := make([][]nodes.NodeOptions, populations)
	candidatesLoss := make([]float64, populations)

	r := rand.New(rand.NewSource(420))
	params := randomOptions(initialParams, r)

	pb := utils.NewDefaultProgressBar("Generations", generations, 50)
	tuner := RandomGeneticEvolve{params, candidates, candidatesLoss, 0, 1, populations, generations, r, pb}
	tuner.generateNewCadidates(initialParams)

	return &tuner
}

func (rs *RandomGeneticEvolve) bestCandidate() []nodes.NodeOptions {
	var index int
	var bestLoss float64 = math.Inf(1)
	for i, loss := range rs.candidatesLoss {
		if loss < bestLoss {
			index = i
			bestLoss = loss
		}
	}
	return rs.candidates[index]
}

func (rs *RandomGeneticEvolve) mutationFactor() float32 {
	var sigma0 float32 = 0.75
	var sigmaMin float32 = 0.0001
	tua := float64(rs.generations) / 2.4
	s := sigma0 * float32(math.Exp(-float64(rs.generation)/tua))
	return max(s, sigmaMin) * float32(rs.r.NormFloat64())
}

func (rs *RandomGeneticEvolve) generateNewCadidates(parent []nodes.NodeOptions) {
	for i := range rs.populations {
		if rs.candidates[i] == nil {
			rs.candidates[i] = make([]nodes.NodeOptions, len(parent))
		}

		var m float32

		for j, opts := range parent {
			params := make([]nodes.NodeParam, 0)

			enabled := opts.Enabled
			if rs.r.Float32() < float32(math.Abs(float64(rs.mutationFactor()))) {
				enabled = !enabled
			}

			for _, param := range opts.Params {

				pEnabled := param.Enabled
				if rs.r.Float32() < float32(math.Abs(float64(rs.mutationFactor()))) {
					pEnabled = !pEnabled
				}

				m = rs.mutationFactor()
				v := min(max(param.Value+m, -1), 1)

				params = append(params, nodes.NodeParam{
					Enabled: pEnabled,
					Name:    param.Name,
					Value:   v,
				})
			}

			rs.candidates[i][j] = nodes.NodeOptions{
				Name:    opts.Name,
				Enabled: enabled,
				Params:  params,
			}
		}
	}
}

func (rs *RandomGeneticEvolve) Update(loss float64) bool {
	rs.candidatesLoss[rs.iteration] = loss

	if rs.iteration < len(rs.candidates)-1 {
		rs.iteration++
		return false
	}

	if rs.generation == rs.generations {
		for i, nodeOpts := range rs.bestCandidate() {
			rs.nodeOpts[i] = nodes.NodeOptions{
				Name:    nodeOpts.Name,
				Enabled: nodeOpts.Enabled,
				Params:  nodeOpts.Params,
			}
		}
		rs.pb.Step()
		return true
	}

	rs.iteration = 0
	rs.generation++
	rs.pb.Step()

	best := rs.bestCandidate()
	rs.generateNewCadidates(best)

	return false
}

func (rs *RandomGeneticEvolve) NodeOptions() []nodes.NodeOptions {
	return rs.nodeOpts
}

func (rs *RandomGeneticEvolve) Candidate() []nodes.NodeOptions {
	return rs.candidates[rs.iteration]
}
