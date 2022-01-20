package adaptive

import (
	"gmidarch/development/messages"
	"shared"
)

//@Type: Planner
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Planner struct {}

func (Planner) I_Process (id string, msg *messages.SAMessage, info *interface{}) {
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
