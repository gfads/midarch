package adaptive

import (
	"gmidarch/development/messages"
	"shared"
	"shared/pluginUtils"
)

//@Type: Executor
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Executor struct {}

func (Executor) I_Process(id string, msg *messages.SAMessage, info *interface{}) {
	plan := msg.Payload.(shared.AdaptationPlan)
	unitCommand := shared.UnitCommand{}

	if len(plan.Operations) > 0 { // TODO
		pluginName := plan.Params[plan.Operations[0]][0]
		plg := pluginUtils.LoadPlugin(pluginName)
		getType, _ := plg.Lookup("Gettype")
		elemType := getType.(func() interface{})()
		getSelector, _ := plg.Lookup("Getselector")
		funcSelector := getSelector.(func() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool))()

		unitCommand.Cmd = shared.REPLACE_COMPONENT
		unitCommand.Params = plg
		unitCommand.Type = elemType
		unitCommand.Selector = funcSelector
	}
	*msg = messages.SAMessage{Payload: unitCommand}
	//fmt.Printf("Executor:: %v\n",msg.Payload)
}