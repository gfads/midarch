package components

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	messages2 "gmidarch/development/messages"
	shared2 "shared"
)

type Analyser struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewAnalyser() Analyser {

	// create a new instance of Server
	r := new(Analyser)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (Analyser) I_Process(msg *messages2.SAMessage, info [] *interface{}) {
	monitoredEvolutiveData := msg.Payload.(shared2.MonitoredEvolutiveData)
	evolutiveAnalysisResult := shared2.EvolutiveAnalysisResult{}

	if len(monitoredEvolutiveData) > 0 { // New plugins available
		evolutiveAnalysisResult.NeedAdaptation = true
		evolutiveAnalysisResult.MonitoredEvolutiveData = monitoredEvolutiveData
	} else {
		evolutiveAnalysisResult.NeedAdaptation = false
		evolutiveAnalysisResult.MonitoredEvolutiveData = monitoredEvolutiveData
	}
	*msg = messages2.SAMessage{Payload: evolutiveAnalysisResult}
}
