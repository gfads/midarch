package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/shared"
)

type Analyser struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewAnalyser() Analyser {

	// create a new instance of Server
	r := new(Analyser)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

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
