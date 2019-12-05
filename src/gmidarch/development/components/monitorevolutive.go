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

func (e Monevolutive) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	e.I_Checkplugins(msg, info)
}

func (Monevolutive) I_Checkplugins(msg *messages.SAMessage, info [] *interface{}) {
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if isFirstTime {
		//time.Sleep(10 * time.Second)  // TODO - only first time
		time.Sleep(1 * time.Second)  // TODO - only first time
		isFirstTime = false
		listOfOldPlugins = shared.LoadPlugins()
	} else {
		time.Sleep(shared.MONITOR_TIME) // TODO
		listOfNewPlugins = shared.LoadPlugins()
		newPlugins = shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
	evolutiveMonitoredData = newPlugins
	*msg = messages.SAMessage{evolutiveMonitoredData}

	listOfOldPlugins = listOfNewPlugins
}
