package managing

import (
	"plugin"
	"sync"
)

type ManagingSystem struct {
}

func (ManagingSystem) Start(fromEnv chan map[string]plugin.Plugin, fromManaged chan int, toManaged chan func(), wg *sync.WaitGroup) {
	toAnalyser := make(chan InfoToAnalyser)
	toPlanner := make(chan InfoToAnalyser)
	toExecutor := make(chan AdaptationPlan)

	go Monitor{}.Start(fromEnv, fromManaged, toAnalyser)
	go Analyser{}.Start(toAnalyser, toPlanner)
	go Planner{}.Start(toPlanner, toExecutor)
	go Executor{}.Start(toExecutor, toManaged)

	wg.Done()
}
