package tuning

type ParamTuner interface {
	Update(loss float64)
	Params() map[string]float32
	Candidate() map[string]float32
}
