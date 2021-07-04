package managing

import (
	"adaptive/adaptiveV4/sharedadaptive"
)

type PlannerInfo struct { // Adaptation plan
	Source  int
	Command string
	Params  interface{}
}

type Planner interface {
	SetMAPEK(MAPEK)
	ToPlanner(AnalyserInfo)
}

type PlannerImpl struct {
	Mapek MAPEK
	Info  PlannerInfo
}

func NewPlanner() Planner {
	var p Planner

	p = &PlannerImpl{}

	return p
}

func (p *PlannerImpl) SetMAPEK(mapek MAPEK) {
	p.Mapek = mapek
}

func (p *PlannerImpl) ToPlanner(info AnalyserInfo) {

	// Define adaptation plan
	p.Info = PlannerInfo{Source: sharedadaptive.FROM_ENV, Command: sharedadaptive.CMD_UPDATE, Params: info}

	// Send to executor
	p.Mapek.E.ToExecutor(p.Info)
}
