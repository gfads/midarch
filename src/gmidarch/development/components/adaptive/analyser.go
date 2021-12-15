package adaptive

import (
	"gmidarch/development/messages"
	"shared"
)

//@Type: Analyser
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Analyser struct {}

func (Analyser) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	monitoredEvolutiveData := msg.Payload.(shared.MonitoredEvolutiveData)
	evolutiveAnalysisResult := shared.EvolutiveAnalysisResult{}

	if len(monitoredEvolutiveData) > 0 { // New plugins available
		evolutiveAnalysisResult.NeedAdaptation = true
		evolutiveAnalysisResult.MonitoredEvolutiveData = monitoredEvolutiveData
	} else {
		evolutiveAnalysisResult.NeedAdaptation = false
		evolutiveAnalysisResult.MonitoredEvolutiveData = monitoredEvolutiveData
	}
	*msg = messages.SAMessage{Payload: evolutiveAnalysisResult}
}
