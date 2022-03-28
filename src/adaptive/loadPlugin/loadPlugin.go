package main

import (
	"fmt"
	"gmidarch/development/messages"
	"gmidarch/development/repositories/architectural"
	"injector"
	"shared"
	"shared/pluginUtils"
)

func main() {
	newPlugins := []string{}
	newPlugins = append(newPlugins, "sender_v1")

	plan := shared.AdaptationPlan{}
	plan.Operations = []string{}
	plan.Params = make(map[string][]string)

	if true {// evolutiveAnalysisResult.NeedAdaptation { // Adaptation is necessary // TODO
		plan.Operations = append(plan.Operations, shared.REPLACE_COMPONENT)
		plan.Params[plan.Operations[0]] = newPlugins
	}

	var msg = &messages.SAMessage{Payload: plan}
	//*msg = messages.SAMessage{Payload: plan}

	pluginName := plan.Params[plan.Operations[0]][0]
	evolutive.GeneratePlugin(pluginName, pluginName)

	I_Process("", msg, nil)
}

func I_Process(id string, msg *messages.SAMessage, info *interface{}) {
	fmt.Println("Executor::msg.Payload", msg.Payload)
	plan := msg.Payload.(shared.AdaptationPlan)
	fmt.Println("Executor::plan", plan)
	unitCommand := shared.UnitCommand{Cmd: "Nothing"}

	if len(plan.Operations) > 0 { // TODO
		pluginName := plan.Params[plan.Operations[0]][0]
		fmt.Println("Executor.I_Process::will load plugin:", pluginName)
		plg := pluginUtils.LoadPlugin(pluginName)
		fmt.Println("Executor.I_Process::plugin loaded:", pluginName)
		fmt.Println("Executor.I_Process::will lookup Gettype:", pluginName)
		getType, _ := plg.Lookup("Gettype")
		//getType := getPluginType(pluginName)
		fmt.Println("Executor.I_Process::loaded Gettype")
		elemType := getType.(func() interface{})()

		fmt.Println("Executor.I_Process::will lookup: Getselector")
		getSelector, _ := plg.Lookup("Getselector")
		funcSelector := getSelector.(func() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool))()

		unitCommand.Cmd = shared.REPLACE_COMPONENT
		unitCommand.Params = plg
		unitCommand.Type = elemType
		unitCommand.Selector = funcSelector
	}
	*msg = messages.SAMessage{Payload: unitCommand}
	fmt.Println("Executor::msg.Payload", msg.Payload)
}

func getPluginType(pluginName string) string {
	pluginSource := shared.DIR_PLUGINS_SOURCE + "/" + pluginName + "/" + pluginName + ".go"
	pluginType, _ := architectural.GetTypeAndBehaviour(pluginSource)
	return pluginType
}