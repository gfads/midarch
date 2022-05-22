package adaptive

import (
	"fmt"
	"gmidarch/development/messages"
	"shared"
	"shared/pluginUtils"
	"time"
)

var isFirstTime = true
var listOfOldPlugins map[string]time.Time

//@Type: Monevolutive
//@Behaviour: Behaviour = I_Hasnewplugins -> InvR.e1 -> Behaviour
type Monevolutive struct {} //[] I_Nonewplugins -> Behaviour

func (Monevolutive) I_Nonewplugins(id string, msg *messages.SAMessage, info *interface{}, reset *bool) { //, r *bool
	listOfNewPlugins := pluginUtils.LoadPlugins()
	newPlugins := pluginUtils.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	if len(newPlugins) != 0 {
		//*r = false
		return
	}
}

func (Monevolutive) I_Hasnewplugins(id string, msg *messages.SAMessage, info *interface{}, reset *bool) { //, r *bool
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if isFirstTime {
		time.Sleep(shared.FIRST_MONITOR_TIME) // only first time
		isFirstTime = false
		listOfOldPlugins = pluginUtils.LoadPlugins()
		fmt.Println("Monevolutive.I_Hasnewplugins::FirstTime - OldPlugins:", listOfOldPlugins)
	} else {
		time.Sleep(shared.MONITOR_TIME)
		listOfNewPlugins = pluginUtils.LoadPlugins()
		newPlugins = pluginUtils.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
		fmt.Println("Monevolutive.I_Hasnewplugins::OldPlugins:", listOfOldPlugins, "NewPlugins:",newPlugins)
	}

	// return from this point if no new pluginsSrc detected
	if len(newPlugins) == 0 {
		fmt.Println("Monevolutive.I_Hasnewplugins::No new pluginsSrc found")
		*reset = true
		//evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
		//evolutiveMonitoredData = newPlugins
		//*msg = messages.SAMessage{Payload: evolutiveMonitoredData}
		return
	}

	fmt.Println("Monevolutive.I_Hasnewplugins::Found new pluginsSrc")
	evolutiveMonitoredData := shared.MonitoredEvolutiveData{} // Todo dcruzb: remove this line, it's overridden by the next one, make some tests on types
	evolutiveMonitoredData = newPlugins
	*msg = messages.SAMessage{Payload: evolutiveMonitoredData}

	listOfOldPlugins = listOfNewPlugins
}
