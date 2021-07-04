package managing

import (
	sharedadaptive2 "adaptive/adaptiveV2/sharedadaptive"
	"plugin"
	"time"
)

type Monitor struct{}

func (Monitor) Start(fromEnv chan map[string]plugin.Plugin, fromManaged chan int, toAnalyser chan InfoToAnalyser) {
	for {
		select {
		case plugins := <-fromEnv:
			toAnalyser <- InfoToAnalyser{Source: sharedadaptive2.FROM_ENV, FromEnv: plugins}
		case n := <-fromManaged:
			toAnalyser <- InfoToAnalyser{Source: sharedadaptive2.FROM_MANAGED, FromManaged: n}
		}
		time.Sleep(sharedadaptive2.MONITOR_TIME * time.Second)
	}
}

