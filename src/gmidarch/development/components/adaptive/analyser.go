package adaptive

import (
	"github.com/gfads/midarch/src/gmidarch/development/messages"
	"github.com/gfads/midarch/src/shared"
)

// @Type: Analyser
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Analyser struct{}

func (Analyser) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	monitoredEvolutiveData := msg.Payload.(shared.MonitoredEvolutiveData)
	evolutiveAnalysisResult := shared.EvolutiveAnalysisResult{}

	//fmt.Println("Analyser.I_Process::monitoredEvolutiveData:", monitoredEvolutiveData)
	if len(monitoredEvolutiveData) > 0 { // New pluginsSrc available
		evolutiveAnalysisResult.NeedAdaptation = true
		evolutiveAnalysisResult.MonitoredEvolutiveData = monitoredEvolutiveData
	} else {
		evolutiveAnalysisResult.NeedAdaptation = false
		evolutiveAnalysisResult.MonitoredEvolutiveData = monitoredEvolutiveData
	}
	*msg = messages.SAMessage{Payload: evolutiveAnalysisResult}
	//fmt.Println("Analyser.I_Process::evolutiveAnalysisResult:", evolutiveAnalysisResult)
}
