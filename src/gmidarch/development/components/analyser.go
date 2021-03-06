package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
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

func (e Analyser) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg, info)
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
