package tuning

import (
	"fmt"
	"maps"
	"math"
	"math/rand"
)

func randomParams(params map[string]float32, r *rand.Rand) map[string]float32 {
	candidate := make(map[string]float32)
	for k := range params {
		candidate[k] = r.Float32()
	}
	return candidate
}

type RandomGeneticEvolve struct {
	params map[string]float32

	candidates     []map[string]float32
	candidatesLoss []float64
	iteration      int
	generation     int

	populations int
	generations int

	r *rand.Rand
}

func NewRandomGeneticEvolve(initialParams map[string]float32, populations, generations int) *RandomGeneticEvolve {
	candidates := make([]map[string]float32, populations)
	candidatesLoss := make([]float64, populations)

	r := rand.New(rand.NewSource(420))
	params := randomParams(initialParams, r)

	tuner := RandomGeneticEvolve{params, candidates, candidatesLoss, 0, 1, populations, generations, r}
	tuner.generateNewCadidates(initialParams)

	return &tuner
}

func (rs *RandomGeneticEvolve) bestCandidate() map[string]float32 {
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
	var sigmaMin float32 = 0.005
	tua := float64(rs.generations) / 2.4
	s := sigma0 * float32(math.Exp(-float64(rs.generation)/tua))
	return max(s, sigmaMin) * float32(rs.r.NormFloat64())
}

func (rs *RandomGeneticEvolve) generateNewCadidates(parent map[string]float32) {
	for i := range rs.populations {
		if rs.candidates[i] == nil {
			rs.candidates[i] = make(map[string]float32)
		}

		for k, v := range parent {
			m := rs.mutationFactor()
			rs.candidates[i][k] = min(max(v+m, -1), 1)
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
		maps.Copy(rs.params, rs.bestCandidate())
		return true
	}

	rs.iteration = 0
	rs.generation++
	fmt.Printf("Starting generation %d\n", rs.generation)

	best := rs.bestCandidate()
	rs.generateNewCadidates(best)

	return false
}

func (rs *RandomGeneticEvolve) Params() map[string]float32 {
	return rs.params
}

func (rs *RandomGeneticEvolve) Candidate() map[string]float32 {
	return rs.candidates[rs.iteration]
}
