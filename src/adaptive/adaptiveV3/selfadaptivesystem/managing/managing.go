package managing

import (
	"plugin"
	"sync"
)

type ManagingSystem struct {}

func (ManagingSystem) StartLocal(fromEnv chan map[string]plugin.Plugin, fromManaged chan int, toManaged chan func(), wg *sync.WaitGroup) {
	toAnalyser := make(chan InfoToAnalyser)
	toPlanner := make(chan InfoToAnalyser)
	toExecutor := make(chan AdaptationPlan)

	go Analyser{}.StartRemote(toAnalyser, toPlanner) // receive from a remote monitor
	go Planner{}.Start(toPlanner, toExecutor)
	go Executor{}.Start(toExecutor, toManaged)

	wg.Done()
}

func (ManagingSystem) StartRemote(fromEnv chan map[string]string, wg *sync.WaitGroup) {

	go Monitor{}.StartRemote(fromEnv)

	wg.Done()
}