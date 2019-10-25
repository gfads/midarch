package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
	"time"
)

var isFirstTime = true
var listOfOldPlugins map[string]time.Time

type Monevolutive struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewMonevolutive() Monevolutive {

	r := new(Monevolutive)
	r.Behaviour = "B = I_Collect -> InvR.e1 -> B"

	return *r
}

func (Monevolutive) I_Checkplugins(msg *messages.SAMessage, info [] *interface{}) {
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if isFirstTime {
		isFirstTime = false
		listOfOldPlugins = shared.LoadPlugins()
	} else {
		listOfNewPlugins = shared.LoadPlugins()
		newPlugins = shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
	evolutiveMonitoredData = newPlugins
	*msg = messages.SAMessage{evolutiveMonitoredData}

	listOfOldPlugins = listOfNewPlugins

	time.Sleep(1000 * time.Millisecond)
}
