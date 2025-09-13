package tuning

import (
	"maps"
	"math"
)

type CoordinateSearch struct {
	LearningRate float32 // initial step
	MinStep      float32 // stop changing when step <= MinStep

	params    map[string]float32
	keys      []string
	bestLoss  float64
	index     int // current coordinate index
	phase     int // 0=baseline, 1=+step, 2=-step
	base      map[string]float32
	candidate map[string]float32
	step      float32

	// sweep bookkeeping
	improvedInSweep bool
	visited         int
}

func NewCoordinateSearch(lr float32, initialParams map[string]float32) *CoordinateSearch {
	keys := make([]string, 0, len(initialParams))
	for k := range initialParams {
		keys = append(keys, k)
	}
	p := maps.Clone(initialParams)
	return &CoordinateSearch{
		LearningRate: lr,
		MinStep:      1e-6, // you can expose/tune this if desired
		params:       p,
		keys:         keys,
		bestLoss:     math.Inf(1),
		index:        0,
		phase:        0,
		base:         maps.Clone(p),
		candidate:    maps.Clone(p), // start by evaluating baseline
		step:         lr,
	}
}

func (cs *CoordinateSearch) Params() map[string]float32 {
	return cs.params
}

func (cs *CoordinateSearch) Candidate() map[string]float32 {
	return cs.candidate
}

// Update consumes the loss of Candidate() and prepares the next candidate.
func (cs *CoordinateSearch) Update(loss float64) {
	// If step is too small, just keep returning current params (converged).
	if cs.step <= cs.MinStep {
		cs.candidate = maps.Clone(cs.params)
		cs.phase = 0
		cs.bestLoss = loss // keep it current
		return
	}

	k := cs.keys[cs.index]

	switch cs.phase {
	case 0: // baseline evaluated
		cs.bestLoss = loss
		cs.base = maps.Clone(cs.params)
		cs.candidate = maps.Clone(cs.base)
		cs.candidate[k] = cs.base[k] + cs.step
		cs.phase = 1
		return

	case 1: // +step evaluated
		if loss < cs.bestLoss {
			// accept +step; continue stepping + on same coord
			cs.bestLoss = loss
			maps.Copy(cs.params, cs.candidate)
			cs.base = maps.Clone(cs.params)
			cs.candidate = maps.Clone(cs.base)
			cs.candidate[k] = cs.base[k] + cs.step
			cs.improvedInSweep = true
			// stay in phase 1
			return
		}
		// try -step from the same base (not from +step)
		cs.candidate = maps.Clone(cs.base)
		cs.candidate[k] = cs.base[k] - cs.step
		cs.phase = 2
		return

	case 2: // -step evaluated
		if loss < cs.bestLoss {
			// accept -step; continue stepping - on same coord
			cs.bestLoss = loss
			maps.Copy(cs.params, cs.candidate)
			cs.base = maps.Clone(cs.params)
			cs.candidate = maps.Clone(cs.base)
			cs.candidate[k] = cs.base[k] - cs.step
			cs.improvedInSweep = true
			// stay in phase 2
			return
		}

		// Neither direction helped; move to next coordinate
		cs.index = (cs.index + 1) % len(cs.keys)
		cs.visited++

		// Completed a sweep without any improvement? Shrink step.
		if cs.visited >= len(cs.keys) {
			if !cs.improvedInSweep {
				cs.step *= 0.5
			}
			cs.visited = 0
			cs.improvedInSweep = false
		}

		// Prepare next coordinate trial
		k = cs.keys[cs.index]
		cs.base = maps.Clone(cs.params)
		cs.candidate = maps.Clone(cs.base)
		// If step got too small after shrink, stick to baseline
		if cs.step <= cs.MinStep {
			cs.phase = 0
			return
		}
		cs.candidate[k] = cs.base[k] + cs.step
		cs.phase = 1
		return
	}
}
