package managing

import (
	"adaptive/adaptiveV2/sharedadaptive"
	"plugin"
	"shared"
)

type Executor struct{}

func (Executor) Start(fromPlanner chan AdaptationPlan, toManaged chan func()) {
	for {
		plan := <-fromPlanner
		switch plan.Command {
		case sharedadaptive.CMD_UPDATE:
			if plan.Source == sharedadaptive.FROM_ENV {
				var f plugin.Symbol
				plugins := plan.Params.(map[string]plugin.Plugin)
				last := ""
				for i := range plugins {
					if i > last {   // Take the most recente plugin only
						last = i
						p := plugins[i]
						f, _ = p.Lookup("Behaviour")
					}
				}
				toManaged <- f.(func())
			} else {
				shared.ErrorHandler(shared.GetFunction(), "TODO")
			}
		default:
			shared.ErrorHandler(shared.GetFunction(), "TODO")
		}
	}
}
