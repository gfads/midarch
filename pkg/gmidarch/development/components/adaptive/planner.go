package adaptive

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
)

// @Type: Planner
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Planner struct{}

func (Planner) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("Planner::msg.Payload", msg.Payload)
	evolutiveAnalysisResult := msg.Payload.(shared.EvolutiveAnalysisResult)
	//fmt.Println("Planner::evolutiveAnalysisResult", evolutiveAnalysisResult)

	plan := shared.AdaptationPlan{}
	plan.Operations = []string{}
	plan.Params = make(map[string][]string)

	if evolutiveAnalysisResult.NeedAdaptation { // Adaptation is necessary // TODO
		plan.Operations = append(plan.Operations, shared.REPLACE_COMPONENT)
		plan.Params[plan.Operations[0]] = evolutiveAnalysisResult.MonitoredEvolutiveData
	}

	*msg = messages.SAMessage{Payload: plan}
	//fmt.Printf("Planner::msg.Payload %v\n",msg.Payload.(shared.AdaptationPlan))
}
