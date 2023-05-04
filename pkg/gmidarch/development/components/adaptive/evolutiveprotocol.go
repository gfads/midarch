package adaptive

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
	"time"
)

var isFirstTimeEvolutiveProtocol = true

//var listOfOldPluginsEvolutiveProtocol map[string]time.Time

// @Type: EvolutiveProtocol
// @Behaviour: Behaviour = I_Hasnewprotocol -> InvR.e1 -> Behaviour
type EvolutiveProtocol struct{} //[] I_Nonewplugins -> Behaviour

//func (EvolutiveProtocol) I_Nonewplugins(id string, msg *messages.SAMessage, info *interface{}, reset *bool) { //, r *bool
//	listOfNewPlugins := pluginUtils.LoadPlugins()
//	newPlugins := pluginUtils.CheckForNewPlugins(listOfOldPluginsEvolutiveProtocol, listOfNewPlugins)
//	if len(newPlugins) != 0 {
//		//*r = false
//		return
//	}
//}

func (EvolutiveProtocol) I_Hasnewprotocol(id string, msg *messages.SAMessage, info *interface{}, reset *bool) { //, r *bool
	//newPlugins := []string{}
	//listOfNewPlugins := make(map[string]time.Time)

	if isFirstTimeEvolutiveProtocol {
		time.Sleep(shared.FIRST_MONITOR_TIME) // only first time
		isFirstTimeEvolutiveProtocol = false
		//listOfOldPluginsEvolutiveProtocol = pluginUtils.LoadPlugins()
		//fmt.Println("EvolutiveProtocol.I_Hasnewplugins::FirstTime - OldPlugins:", listOfOldPluginsEvolutiveProtocol)
	} else {
		time.Sleep(shared.MONITOR_TIME)
		//listOfNewPlugins = pluginUtils.LoadPlugins()
		//newPlugins = pluginUtils.CheckForNewPlugins(listOfOldPluginsEvolutiveProtocol, listOfNewPlugins)
		//fmt.Println("EvolutiveProtocol.I_Hasnewplugins::OldPlugins:", listOfOldPluginsEvolutiveProtocol, "NewPlugins:",newPlugins)
	}

	// return from this point if no new pluginsSrc detected
	if len(shared.ListOfComponentsToAdaptTo) == 0 {
		//fmt.Println("EvolutiveProtocol.I_Hasnewplugins::No new pluginsSrc found")
		*reset = true
		//evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
		//evolutiveMonitoredData = newPlugins
		//*msg = messages.SAMessage{Payload: evolutiveMonitoredData}
		return
	}

	//fmt.Println("EvolutiveProtocol.I_Hasnewplugins::Found new pluginsSrc")
	evolutiveMonitoredData := shared.MonitoredEvolutiveData{} // Todo dcruzb: remove this line, it's overridden by the next one, make some tests on types
	evolutiveMonitoredData = shared.ListOfComponentsToAdaptTo
	*msg = messages.SAMessage{Payload: evolutiveMonitoredData}

	//listOfOldPluginsEvolutiveProtocol = listOfNewPlugins
	shared.ListOfComponentsToAdaptTo = []string{}
}
