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
	//r.Behaviour = "B = I_Checkplugins -> InvR.e1 -> B"
	r.Behaviour = "B = I_Hasnewplugins -> InvR.e1 -> B [] I_Nonewplugins -> B"

	return *r
}

func (e Monevolutive) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'H' {
		e.I_Hasnewplugins(msg, info, r)
	} else {
		e.I_Nonewplugins(msg, info, r)
	}
}

func (Monevolutive) I_Nonewplugins(msg *messages.SAMessage, info [] *interface{}, r *bool) {
	listOfNewPlugins := shared.LoadPlugins()
	newPlugins := shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	if len(newPlugins) != 0 {
		*r = false
		return
	}
}

func (Monevolutive) I_Hasnewplugins(msg *messages.SAMessage, info [] *interface{}, r *bool) {
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if isFirstTime {
		time.Sleep(shared.FIRST_MONITOR_TIME) // only first time
		isFirstTime = false
		listOfOldPlugins = shared.LoadPlugins()
	} else {
		time.Sleep(shared.MONITOR_TIME)
		listOfNewPlugins = shared.LoadPlugins()
		newPlugins = shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	// return from this point if no new plugins detected
	if len(newPlugins) == 0 {
		*r = false
		return
	}

	evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
	evolutiveMonitoredData = newPlugins
	*msg = messages.SAMessage{evolutiveMonitoredData}

	listOfOldPlugins = listOfNewPlugins
}
