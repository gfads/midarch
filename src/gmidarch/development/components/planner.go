package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type Planner struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewPlanner() Planner {

	// create a new instance of Server
	r := new(Planner)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (e Planner) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	e.I_Process(msg, info)
}

func (Planner) I_Process (msg *messages.SAMessage, info [] *interface{}) {
	evolutiveAnalysisResult := msg.Payload.(shared.EvolutiveAnalysisResult)

	plan := shared.AdaptationPlan{}
	plan.Operations = []string{}
	plan.Params = make(map[string][]string)

	if evolutiveAnalysisResult.NeedAdaptation { // Adaptation is necessary // TODO
		plan.Operations = append(plan.Operations, shared.REPLACE_COMPONENT)
		plan.Params[plan.Operations[0]] = evolutiveAnalysisResult.MonitoredEvolutiveData
	}

	*msg = messages.SAMessage{Payload: plan}
}
