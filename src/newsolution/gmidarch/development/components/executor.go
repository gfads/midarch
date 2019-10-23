package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/shared"
)

type Executor struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewExecutor() Executor {

	// create a new instance of client
	r := new(Executor)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (Executor) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	plan := msg.Payload.(shared.AdaptationPlan)

	unitCommand := shared.UnitCommand{}

	if len(plan.Operations) > 0 { // TODO
	    pluginName := plan.Params[plan.Operations[0]][0]
		plg := shared.LoadPlugin(pluginName)
		tp,_ := plg.Lookup("Gettype")

		unitCommand.Cmd = "STOP"
		unitCommand.Params = plg
		unitCommand.Type = tp
	}
	*msg = messages.SAMessage{Payload:unitCommand}
}
