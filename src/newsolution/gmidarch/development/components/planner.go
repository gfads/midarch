package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
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

func (Planner) I_Createplan (msg *messages.SAMessage, info [] *interface{}) {
	evolutiveAnalysisResult := msg.Payload.(shared.EvolutiveAnalysisResult)

	plan := shared.AdaptationPlan{}
	plan.Operations = []string{}
	plan.Params = make(map[string][]string)

	if evolutiveAnalysisResult.NeedAdaptation { // Adaptation is necessary // TODO
		plan.Operations = append(plan.Operations, parameters.REPLACE_COMPONENT)
		plan.Params[plan.Operations[0]] = evolutiveAnalysisResult.MonitoredEvolutiveData
	}

	*msg = messages.SAMessage{Payload: plan}
}
