package components

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/shared"
	"time"
)

var isFirstTime = true
var listOfOldPlugins map[string]time.Time

type Monevolutive struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewMonevolutive() Monevolutive {

	// create a new instance of Server
	r := new(Monevolutive)
	r.Behaviour = "B = I_Collect -> InvR.e1 -> B"

	return *r
}

func (Monevolutive) I_Checkplugins(msg *messages.SAMessage, info [] *interface{}) {
	//confName := (*info).(string)
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if isFirstTime {
		isFirstTime = false
		//listOfOldPlugins = shared.LoadPlugins(confName)
		listOfOldPlugins = shared.LoadPlugins()
	} else {
		//listOfNewPlugins = shared.LoadPlugins(confName)
		listOfNewPlugins = shared.LoadPlugins()
		newPlugins = shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	if len(newPlugins) > 0 {
		evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
		evolutiveMonitoredData = newPlugins
		*msg = messages.SAMessage{evolutiveMonitoredData}
	}

	listOfOldPlugins = listOfNewPlugins
	//time.Sleep(parameters.MONITOR_TIME * time.Second)
	time.Sleep(1 * time.Second)
	fmt.Printf("Monitoevolutive:: EXIT %v\n",*msg)
}