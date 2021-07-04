package managing

type MAPEK struct {
	M Monitor
	A Analyser
	P Planner
	E Executor
}

func NewMAPEK(m Monitor, a Analyser, p Planner, e Executor) MAPEK {
	r := MAPEK{}

	r.M = m
	r.A = a
	r.P = p
	r.E = e

	return r
}
