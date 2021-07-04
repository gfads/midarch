package main

import (
	"adaptive/adaptiveV4/environment/plugins/injector"
	"adaptive/adaptiveV4/selfadaptivesystem/managed"
	"adaptive/adaptiveV4/selfadaptivesystem/managing"
	"sync")

func main() {
	var wg sync.WaitGroup

	// Configure probes of the monitor
	//var probes []func(managing.MonitorImpl) managing.MonitorInfo
	//probes = append(probes, managing.MonitorImpl.ProbeSourceRemote)
	//m := managing.NewMonitor(probes) // no managed systems

	// Configure MAPE-K (no monitor)
	mapek := managing.NewMAPEK(nil, managing.NewAnalyser(), managing.NewPlanner(), managing.NewExecutor())
	managedSystem := managed.NewManaged()

	// Configure managing system
	managingSystem := managing.NewManagingSystem(managedSystem, &mapek)

	// Empty repositories
	inj := injector.PluginInjector{}
	inj.Initialize()

	// Start elements
	wg.Add(2)
	go managedSystem.Start(&wg)
	go managingSystem.Start(&wg)

	//go AdaptationGoals()

	wg.Wait()
}
