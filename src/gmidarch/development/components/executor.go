package components

import (
	"fmt"
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

func (e Executor) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	e.I_Process(msg, info)
}

func (Executor) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	plan := msg.Payload.(shared.AdaptationPlan)
	unitCommand := shared.UnitCommand{}

	if len(plan.Operations) > 0 { // TODO
		pluginName := plan.Params[plan.Operations[0]][0]
		plg := shared.LoadPlugin(pluginName)
		getType, _ := plg.Lookup("Gettype")
		elemType := getType.(func() interface{})()
		getSelector, _ := plg.Lookup("Getselector")
		funcSelector := getSelector.(func() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}))()

		unitCommand.Cmd = shared.REPLACE_COMPONENT
		unitCommand.Params = plg
		unitCommand.Type = elemType
		unitCommand.Selector = funcSelector
	}
	*msg = messages.SAMessage{Payload: unitCommand}
	fmt.Printf("Executor:: %v\n",msg.Payload)
}
