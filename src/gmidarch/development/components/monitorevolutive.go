package components

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	messages2 "gmidarch/development/messages"
	shared2 "shared"
	"time"
)

var isFirstTime = true
var listOfOldPlugins map[string]time.Time

type Monevolutive struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewMonevolutive() Monevolutive {

	r := new(Monevolutive)
	r.Behaviour = "B = I_Collect -> InvR.e1 -> B"

	return *r
}

func (Monevolutive) I_Checkplugins(msg *messages2.SAMessage, info [] *interface{}) {
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if isFirstTime {
		isFirstTime = false
		listOfOldPlugins = shared2.LoadPlugins()
	} else {
		listOfNewPlugins = shared2.LoadPlugins()
		newPlugins = shared2.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	evolutiveMonitoredData := shared2.MonitoredEvolutiveData{}
	evolutiveMonitoredData = newPlugins
	*msg = messages2.SAMessage{evolutiveMonitoredData}

	listOfOldPlugins = listOfNewPlugins

	time.Sleep(1000 * time.Millisecond)
}
