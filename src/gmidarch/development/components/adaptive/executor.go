package adaptive

import (
	"fmt"
	"gmidarch/development/messages"
	"log"
	"reflect"
	"shared"
	"shared/pluginUtils"
)

//@Type: Executor
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Executor struct{}

func (Executor) I_Process(id string, msg *messages.SAMessage, info *interface{}) {
	fmt.Println("Executor::msg.Payload", msg.Payload)
	plan := msg.Payload.(shared.AdaptationPlan)
	fmt.Println("Executor::plan", plan)
	unitCommand := shared.UnitCommand{Cmd: "Nothing"}

	if len(plan.Operations) > 0 { // TODO
		pluginName := plan.Params[plan.Operations[0]][0]
		fmt.Println("Executor.I_Process::will load plugin:", pluginName)
		plg := pluginUtils.LoadPlugin(pluginName)
		fmt.Println("Executor.I_Process::plugin loaded:", pluginName)
		log.Println("Executor.I_Process::Will lookup Gettype:", pluginName)
		getType, _ := plg.Lookup("GetType")
		elemType := getType.(func() interface{})()

		log.Println("--------------Executor Adapt to ---->", reflect.TypeOf(elemType).Name())

		//fmt.Println("Executor.I_Process::will lookup Getselector:", pluginName)
		//getSelector, _ := plg.Lookup("Getselector")
		//funcSelector := getSelector.(func() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool))()

		unitCommand.Cmd = shared.REPLACE_COMPONENT
		unitCommand.Params = plg
		unitCommand.Type = elemType
		//unitCommand.Selector = funcSelector
	}
	*msg = messages.SAMessage{Payload: unitCommand}
	fmt.Println("Executor::msg.Payload", msg.Payload)
}
