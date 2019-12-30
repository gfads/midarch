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
	r.Behaviour = "B = I_Checkplugins -> InvR.e1 -> B"
	//r.Behaviour = "B = I_Checkplugins -> InvR.e1 -> B [] I_Noplugins -> B"

	return *r
}

func (e Monevolutive) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	//if op[2] == 'C' {
		e.I_Checkplugins(msg, info, r)
	//} else {
	//	e.I_Noplugins(msg, info, r)
	//}
}

func (Monevolutive) I_Noplugins(msg *messages.SAMessage, info [] *interface{}, r *bool) {
	listOfNewPlugins := make(map[string]time.Time)

	if len(listOfNewPlugins) != 0 {
		*r = false
		return
	}
}

func (Monevolutive) I_Checkplugins(msg *messages.SAMessage, info [] *interface{}, r *bool) {
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	//if len(listOfNewPlugins) == 0 {
	//	*r = false
	//	return
	//}

	if isFirstTime {
		time.Sleep(shared.FIRST_MONITOR_TIME) // TODO - only first time
		//time.Sleep(1 * time.Second)
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
