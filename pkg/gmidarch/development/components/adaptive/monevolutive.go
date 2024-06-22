package adaptive

import (
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/pluginUtils"
)

var isFirstTime = true
var listOfOldPlugins map[string]time.Time

// @Type: Monevolutive
// @Behaviour: Behaviour = I_Hasnewplugins -> InvR.e1 -> Behaviour
type Monevolutive struct{}

func (Monevolutive) I_Hasnewplugins(id string, msg *messages.SAMessage, info *interface{}, reset *bool) { //, r *bool
	evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
	listOfNewPlugins := make(map[string]time.Time)

	if isFirstTime {
		time.Sleep(shared.FIRST_MONITOR_TIME) // only first time
		isFirstTime = false
		listOfOldPlugins = pluginUtils.LoadPlugins()
		//fmt.Println("Monevolutive.I_Hasnewplugins::FirstTime - OldPlugins:", listOfOldPlugins)
	} else {
		time.Sleep(shared.MONITOR_TIME)
		listOfNewPlugins = pluginUtils.LoadPlugins()
		evolutiveMonitoredData = pluginUtils.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
		//fmt.Println("Monevolutive.I_Hasnewplugins::OldPlugins:", listOfOldPlugins, "NewPlugins:",newPlugins)
	}

	// return from this point if no new pluginsSrc detected
	if len(evolutiveMonitoredData) == 0 {
		//fmt.Println("Monevolutive.I_Hasnewplugins::No new pluginsSrc found")
		*reset = true
		return
	}

	//fmt.Println("Monevolutive.I_Hasnewplugins::Found new pluginsSrc")
	*msg = messages.SAMessage{Payload: evolutiveMonitoredData}

	listOfOldPlugins = listOfNewPlugins
}
