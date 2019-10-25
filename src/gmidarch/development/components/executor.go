package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
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
		elemType := tp.(func()interface{})()

		unitCommand.Cmd = shared.REPLACE_COMPONENT
		unitCommand.Params = plg
		unitCommand.Type = elemType
	}
	*msg = messages.SAMessage{Payload: unitCommand}
}
