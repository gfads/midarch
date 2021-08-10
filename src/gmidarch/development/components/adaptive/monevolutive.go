package components

import (
	"gmidarch/development/messages"
	"shared"
	"shared/pluginUtils"
	"time"
)

var isFirstTime = true
var listOfOldPlugins map[string]time.Time

//@Type: Monevolutive
//@Behaviour: Behaviour = B = I_Hasnewplugins -> InvR.e1 -> B [] I_Nonewplugins -> B
type Monevolutive struct {}

func (Monevolutive) I_Nonewplugins(msg *messages.SAMessage, info [] *interface{}, r *bool) {
	listOfNewPlugins := pluginUtils.LoadPlugins()
	newPlugins := pluginUtils.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
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
		listOfOldPlugins = pluginUtils.LoadPlugins()
	} else {
		time.Sleep(shared.MONITOR_TIME)
		listOfNewPlugins = pluginUtils.LoadPlugins()
		newPlugins = pluginUtils.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	// return from this point if no new plugins detected
	if len(newPlugins) == 0 {
		*r = false
		return
	}

	evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
	evolutiveMonitoredData = newPlugins
	*msg = messages.SAMessage{Payload: evolutiveMonitoredData}

	listOfOldPlugins = listOfNewPlugins
}
