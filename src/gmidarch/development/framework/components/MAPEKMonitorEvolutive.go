package components

import (
	"time"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"gmidarch/development/framework/messages"
	"fmt"
)

type MAPEKMonitorEvolutive struct{}

var firstTime = true
var listOfOldPlugins map[string]time.Time

func (MAPEKMonitorEvolutive) I_EvolutiveMonitoring(msg *messages.SAMessage, info *interface{}, r *bool) {

	confName := (*info).(string)
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if firstTime {
		firstTime = false
		listOfOldPlugins = shared.LoadPlugins(confName)
	} else {
		listOfNewPlugins = shared.LoadPlugins(confName)
		newPlugins = shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	if len(newPlugins) > 0 {
		evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
		evolutiveMonitoredData = newPlugins
		*msg = messages.SAMessage{evolutiveMonitoredData}
	}

	listOfOldPlugins = listOfNewPlugins
	time.Sleep(parameters.MONITOR_TIME * time.Second)
}

func (MAPEKMonitorEvolutive) I_HasPlugin(msg *messages.SAMessage, info *interface{}, r *bool) {

	//if msg.Payload != nil {
	if len(listOfOldPlugins) >= 1 {
		*msg = messages.SAMessage{listOfOldPlugins}
		*r = true
	}

	fmt.Printf("MAPEKEvolutiveMonitor:: I_HasPlugin %v \n",msg.Payload)

}

func (MAPEKMonitorEvolutive) I_HasNotPlugin(msg *messages.SAMessage, info *interface{}, r *bool) {

	//if msg.Payload == nil {
	if len(listOfOldPlugins) == 0 {
		*r = true
	}
	fmt.Printf("MAPEKEvolutiveMonitor:: I_HasNotPlugin %v \n",msg.Payload)
}
