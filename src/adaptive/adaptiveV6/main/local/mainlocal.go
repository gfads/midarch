package main

import (
	"adaptive/adaptiveV6/environment/plugins/injector"
	"adaptive/adaptiveV6/selfadaptivesystem/managed"
	"adaptive/adaptiveV6/selfadaptivesystem/managing"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Configure MAPE-K (no monitor)
	mapek := managing.NewMAPEK(nil, managing.NewAnalyser(), managing.NewPlanner(), managing.NewExecutor())
	managedSystem := managed.NewManaged()

	// Configure managing system
	managingSystem := managing.NewManagingSystem(managedSystem, &mapek)

	// Empty repositories
	inj := injector.PluginInjector{}
	inj.Initialize()

	// Start managed and managing systems
	wg.Add(2)
	go managedSystem.Start(&wg)
	go managingSystem.Start(&wg)

	//go AdaptationGoals()   // TODO

	wg.Wait()
}
