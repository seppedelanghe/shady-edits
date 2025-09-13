package tuning

type ParamTuner interface {
	Update(loss float64) bool
	Params() map[string]float32
	Candidate() map[string]float32
}
