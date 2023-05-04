package adaptive

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/pluginUtils"
	"plugin"
	"strings"
)

// @Type: Executor
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Executor struct{}

func (Executor) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("Executor::msg.Payload", msg.Payload)
	plan := msg.Payload.(shared.AdaptationPlan)
	//fmt.Println("Executor::plan", plan)
	unitCommand := shared.UnitCommand{Cmd: "Nothing"}
	//fmt.Println("plan", plan, "plan.Operations", plan.Operations, "plan.Params", plan.Params)
	//shared.ErrorHandler(shared.GetFunction(), "Teste")

	if len(plan.Operations) > 0 { // TODO
		pluginName := plan.Params[plan.Operations[0]][0]
		//fmt.Println("Executor.I_Process::will load plugin:", pluginName)

		if shared.Contains(shared.Adaptability, shared.EVOLUTIVE_PROTOCOL_ADAPTATION) { //strings.Contains(pluginName, "crh") {
			//fmt.Println("EVOLUTIVE_PROTOCOL_ADAPTATION no executor:", pluginName)
			unitCommand.Cmd = shared.REPLACE_COMPONENT
			unitCommand.Params = plugin.Plugin{} //plg
			componentName := strings.ToUpper(strings.ReplaceAll(pluginName, ".so", ""))
			componentName, _, _ = strings.Cut(componentName, "_")
			unitCommand.Type = shared.GetComponentTypeByNameFromRAM(componentName)
			//shared.ErrorHandler(shared.GetFunction(), "Teste")
		} else if shared.Contains(shared.Adaptability, shared.EVOLUTIVE_ADAPTATION) &&
			strings.Contains(pluginName, "srh") { // TODO dcruzb: remove test condition
			//fmt.Println("EVOLUTIVE_ADAPTATION no executor")
			plg := pluginUtils.LoadPlugin(pluginName)
			//fmt.Println("Executor.I_Process::plugin loaded:", pluginName)
			//log.Println("Executor.I_Process::Will lookup Gettype:", pluginName)
			getType, _ := plg.Lookup("GetType")
			elemType := getType.(func() interface{})()

			//log.Println("--------------Executor Adapt to ---->", reflect.TypeOf(elemType).Name())

			//fmt.Println("Executor.I_Process::will lookup Getselector:", pluginName)
			//getSelector, _ := plg.Lookup("Getselector")
			//funcSelector := getSelector.(func() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool))()

			unitCommand.Cmd = shared.REPLACE_COMPONENT
			unitCommand.Params = plg
			unitCommand.Type = elemType
		}

		//log.Println("--------------Executor Adapt to ---->", reflect.TypeOf(elemType).Name())
		//
		////fmt.Println("Executor.I_Process::will lookup Getselector:", pluginName)
		////getSelector, _ := plg.Lookup("Getselector")
		////funcSelector := getSelector.(func() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool))()
		//
		//unitCommand.Cmd = shared.REPLACE_COMPONENT
		//unitCommand.Params = plg
		//unitCommand.Type = elemType
		//unitCommand.Selector = funcSelector
	}
	if unitCommand.Cmd == "Nothing" {
		*reset = true
		return
	}
	*msg = messages.SAMessage{Payload: unitCommand}
	//fmt.Println("Executor::msg.Payload", msg.Payload)
}
