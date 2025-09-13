package tuning

import (
	"maps"
	"math"
	"math/rand"
)

type RandomSearch struct {
	params        map[string]float32
	candidate     map[string]float32
	bestLoss      float64
	iteration     int
	maxIterations int

	r *rand.Rand
}

func NewRandomSearch(initialParams map[string]float32, maxIterations int) *RandomSearch {
	r := rand.New(rand.NewSource(420))
	return &RandomSearch{initialParams, nil, math.Inf(1), 0, maxIterations, r}
}

func (rs *RandomSearch) Update(loss float64) bool {
	if loss < rs.bestLoss {
		rs.bestLoss = loss
		maps.Copy(rs.params, rs.candidate)
	}

	rs.iteration++
	rs.candidate = make(map[string]float32)
	for k := range rs.params {
		rs.candidate[k] = rs.r.Float32()
	}

	return rs.iteration >= rs.maxIterations
}

func (rs *RandomSearch) Params() map[string]float32 {
	return rs.params
}

func (rs *RandomSearch) Candidate() map[string]float32 {
	return rs.candidate
}
