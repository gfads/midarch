package main

import (
	"adaptive/adaptiveV6/environment/plugins/injector"
	"adaptive/adaptiveV6/selfadaptivesystem/managing"
	"adaptive/adaptiveV6/sharedadaptive"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Configure probes of the monitor
	var probes []func(managing.MonitorImpl) interface{}
	probes = append(probes, managing.MonitorImpl.ProbePluginsRemote)
	m := managing.NewMonitor(probes) // no managed systems

	// Configure MAPE-K
	mapek := managing.NewMAPEK(m, nil, nil, nil)              // only monitor
	managingSystem := managing.NewManagingSystem(nil, &mapek) // no managed system

	// Empty & initialise repositories
	inj := injector.PluginInjector{}
	inj.Initialize()

	// Start elements
	wg.Add(2)
	go inj.Start(sharedadaptive.REMOTE, &wg)
	go managingSystem.Start(&wg)

	//go AdaptationGoals()

	wg.Wait()
}
