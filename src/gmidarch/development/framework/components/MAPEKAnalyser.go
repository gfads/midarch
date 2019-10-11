package components

import (
	"newsolution/shared/shared"
	"newsolution/shared/parameters"
	"gmidarch/development/framework/messages"
	"fmt"
	"os"
)

type MAPEKAnalyser struct{}

func (MAPEKAnalyser) I_Analyse(msg *messages.SAMessage, info interface{}, r *bool) {

	fmt.Printf("MAPEAnalyser:: I_Analyse \n")
	os.Exit(0)
	
	// Information received from Monitor
	listOfNewPlugins := msg.Payload.(shared.MonitoredEvolutiveData)

	// TODO Analyse the remaining kind of monitored data
	// Analyse the Evolutive monitored data
	analysisResult := shared.AnalysisResult{}
	analysisResult.Result = listOfNewPlugins
	analysisResult.Analysis = parameters.EVOLUTIVE_CHANGE
	*msg = messages.SAMessage{Payload: analysisResult}
}
