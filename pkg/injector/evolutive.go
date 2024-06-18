package evolutive

import (
	"time"

	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/pluginUtils"
)

type EvolutiveInjector struct{}

func (EvolutiveInjector) Start(firstElem, secondElem string, interval time.Duration) {
	// Replacing strategies
	//go noChange()
	//go changeOnce(firstElem, interval)
	//go changeSamePluginSeveralTimes(elem)
	go alternatePlugins(firstElem, secondElem, interval)
}

func (EvolutiveInjector) StartEvolutiveProtocolInjection(firstElem, secondElem string, interval time.Duration) {
	// Replacing strategies
	//go noChange()
	//go changeOnce(firstElem, interval)
	//go changeSamePluginSeveralTimes(elem)
	go alternateProtocol(firstElem, secondElem, interval)
}

func alternateProtocol(firstElem, secondElem string, interval time.Duration) {
	currentPlugin := 1
	for {
		//fmt.Printf("Evolutive:: Next plugin '%v' will be generated in %v !! \n", elemNew, interval)
		time.Sleep(interval)

		switch currentPlugin {
		case 1: // Plugin 01
			currentPlugin = 2
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, firstElem)
			//elemOld = firstElem + "_v2"
			//elemNew = firstElem + "_v2"
			//GeneratePlugin(elemOld, firstElem, elemNew)
		case 2: // Plugin 02
			currentPlugin = 1
			shared.ListOfComponentsToAdaptTo = append(shared.ListOfComponentsToAdaptTo, secondElem)
			//elemOld = secondElem + "_v1"
			//elemNew = secondElem + "_v1"
			//GeneratePlugin(elemOld, secondElem, elemNew)
		}
	}
}

func changeOnce(elem string, interval time.Duration) {
	//removeOldPlugins()
	time.Sleep(interval)
	elemNew := elem + "_v1"
	pluginUtils.GeneratePlugin(elem, elemNew)
}

// func changeSamePluginSeveralTimes(elem string) {
// 	for {
// 		removeOldPlugins()
// 		generatePlugin(elem, elem+"_v1")
// 		time.Sleep(shared.INJECTION_TIME * time.Second)
// 	}
// }

func alternatePlugins(firstElem, secondElem string, interval time.Duration) {
	//removeOldPlugins()

	elemNew := ""

	currentPlugin := 1
	for {
		//fmt.Printf("Evolutive:: Next plugin '%v' will be generated in %v !! \n", elemNew, interval)
		time.Sleep(interval)

		switch currentPlugin {
		case 1: // Plugin 01
			currentPlugin = 2
			elemNew = firstElem + "_v2"
			pluginUtils.GeneratePlugin(firstElem, elemNew)
		case 2: // Plugin 02
			currentPlugin = 1
			elemNew = secondElem + "_v1"
			pluginUtils.GeneratePlugin(secondElem, elemNew)
		}
	}
}
