package tuning

import "shady-edits/pkg/nodes"

type ParamTuner interface {
	Update(loss float64) bool
	NodeOptions() []nodes.NodeOptions
	Candidate() []nodes.NodeOptions
}
